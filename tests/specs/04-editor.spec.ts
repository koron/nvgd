import { expect, test } from '@playwright/test';
import { OpfsPage } from '../pages/OpfsPage';
import { captureDialog } from '../helpers/dialogs';
import {
  clearOPFS,
  readOPFSFile,
  seedBinaryFile,
  seedFile,
} from '../helpers/opfs';

test.describe('D. Simple editor', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/opfs/');
    await clearOPFS(page);
  });

  test('D1: creates a new file from editor inputs', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.saveEditor('note.txt', 'first contents');

    await expect(opfs.row('note.txt')).toBeVisible();
    expect(await readOPFSFile(page, 'note.txt')).toBe('first contents');
    await expect(opfs.editorName).toHaveValue('');
    await expect(opfs.editorBody).toHaveValue('');
  });

  test('D2: saving without a name alerts and does not create a file', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    const dialogPromise = captureDialog(page, 'accept');
    await opfs.editorBody.fill('orphan body');
    await opfs.editorSaveBtn.click();

    const dialog = await dialogPromise;
    expect(dialog.message()).toMatch(/Need file name/i);
    expect(await opfs.rowCount()).toBe(0);
  });

  test('D3: re-saving updates the existing file', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.saveEditor('a.txt', 'v1');
    await expect(opfs.row('a.txt')).toBeVisible();

    await opfs.saveEditor('a.txt', 'v2-updated');
    await expect
      .poll(() => readOPFSFile(page, 'a.txt'))
      .toBe('v2-updated');
  });

  test('D4: Edit action loads name and content into the editor (< 64KiB)', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'pre.txt', 'pre-existing contents');
    await opfs.goto();

    await opfs.clickEdit('pre.txt');

    await expect(opfs.editorName).toHaveValue('pre.txt');
    await expect(opfs.editorBody).toHaveValue('pre-existing contents');
  });

  test('D5: files >= 64KiB have no Edit action', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await seedBinaryFile(page, 'huge.bin', 64 * 1024 + 1);
    await opfs.goto();

    const row = opfs.row('huge.bin');
    await expect(row).toBeVisible();
    // Save as is always present for files; Edit is hidden for large
    // files. The action <a> contains an icon span + the label, so we
    // match against the substring instead of using exact text.
    await expect(opfs.saveAsAction('huge.bin')).toBeVisible();
    await expect(opfs.editAction('huge.bin')).toHaveCount(0);
  });

  test('D6: Clear empties editor inputs', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.editorName.fill('temp.txt');
    await opfs.editorBody.fill('temp body');

    await opfs.editorClearBtn.click();

    await expect(opfs.editorName).toHaveValue('');
    await expect(opfs.editorBody).toHaveValue('');
  });
});
