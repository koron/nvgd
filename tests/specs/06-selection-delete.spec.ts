import { expect, test } from '@playwright/test';
import { OpfsPage } from '../pages/OpfsPage';
import { captureDialog, autoAcceptDialogs } from '../helpers/dialogs';
import {
  clearOPFS,
  existsOPFS,
  seedDirectory,
  seedFile,
} from '../helpers/opfs';

test.describe('F. Selection & bulk delete', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/opfs/');
    await clearOPFS(page);
  });

  test('F1: ticking a row enables Delete and DuckDB buttons', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'a.txt', 'A');
    await opfs.goto();

    await expect(opfs.deleteBtn).toBeDisabled();
    await expect(opfs.duckdbBtn).toBeDisabled();

    await opfs.selectRow('a.txt');

    await expect(opfs.deleteBtn).toBeEnabled();
    await expect(opfs.duckdbBtn).toBeEnabled();
  });

  test('F2: master checkbox selects/deselects every row', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'a.txt', 'A');
    await seedFile(page, 'b.txt', 'B');
    await seedDirectory(page, 'dirX');
    await opfs.goto();

    await opfs.toggleSelectAll.check();
    const checked = await opfs.grid
      .locator('input.selectedFile')
      .evaluateAll((els) => els.every((e) => (e as HTMLInputElement).checked));
    expect(checked).toBe(true);

    await opfs.toggleSelectAll.uncheck();
    const allClear = await opfs.grid
      .locator('input.selectedFile')
      .evaluateAll((els) => els.every((e) => !(e as HTMLInputElement).checked));
    expect(allClear).toBe(true);
  });

  test('F3: partial selection marks the master checkbox indeterminate', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'a.txt', 'A');
    await seedFile(page, 'b.txt', 'B');
    await opfs.goto();

    await opfs.selectRow('a.txt');

    const indeterminate = await opfs.toggleSelectAll.evaluate(
      (el) => (el as HTMLInputElement).indeterminate,
    );
    expect(indeterminate).toBe(true);
  });

  test('F4: deleting selected files removes them after confirm', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'a.txt', 'A');
    await seedFile(page, 'b.txt', 'B');
    await opfs.goto();

    await opfs.selectRow('a.txt');

    const dialogPromise = captureDialog(page, 'accept');
    await opfs.deleteBtn.click();
    const dialog = await dialogPromise;
    expect(dialog.message()).toMatch(/delete the following/i);

    await expect(opfs.row('a.txt')).toHaveCount(0);
    await expect(opfs.row('b.txt')).toBeVisible();
    await expect(opfs.deleteBtn).toBeDisabled();
  });

  test('F5: deleting a directory removes it recursively', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'parent/child.txt', 'nested');
    await opfs.goto();

    await opfs.selectRow('parent/');
    const dialogPromise = captureDialog(page, 'accept');
    await opfs.deleteBtn.click();
    await dialogPromise;

    expect(await existsOPFS(page, 'parent')).toBeNull();
    expect(await existsOPFS(page, 'parent/child.txt')).toBeNull();
  });

  test('F6: dismissing the confirm leaves files intact', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'keep.txt', 'KEEP');
    await opfs.goto();

    await opfs.selectRow('keep.txt');

    const dialogPromise = captureDialog(page, 'dismiss');
    await opfs.deleteBtn.click();
    await dialogPromise;

    await expect(opfs.row('keep.txt')).toBeVisible();
    expect(await existsOPFS(page, 'keep.txt')).toBe('file');
  });

  test('F7: after deletion the listing stays consistent across reload', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'a.txt', 'A');
    await seedFile(page, 'b.txt', 'B');
    await opfs.goto();

    const detach = autoAcceptDialogs(page);
    await opfs.selectRow('a.txt');
    await opfs.deleteBtn.click();
    detach();

    await expect(opfs.row('a.txt')).toHaveCount(0);
    await opfs.reloadBtn.click();
    await expect(opfs.row('a.txt')).toHaveCount(0);
    await expect(opfs.row('b.txt')).toBeVisible();
  });
});
