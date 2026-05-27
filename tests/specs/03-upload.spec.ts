import { join } from 'node:path';
import { expect, test } from '@playwright/test';
import { OpfsPage } from '../pages/OpfsPage';
import { captureDialog, recordDialogs } from '../helpers/dialogs';
import { clearOPFS, readOPFSFile, statOPFSFile } from '../helpers/opfs';

const FIXTURES = join(__dirname, '..', 'fixtures');
const SMALL = join(FIXTURES, 'small.txt');

test.describe('C. File upload', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/opfs/');
    await clearOPFS(page);
  });

  test('C1: selecting a file auto-fills name and enables Upload', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await expect(opfs.uploadBtn).toBeDisabled();
    await opfs.uploadFileInput.setInputFiles(SMALL);

    await expect(opfs.uploadName).toHaveValue('small.txt');
    await expect(opfs.uploadBtn).toBeEnabled();
  });

  test('C2: upload succeeds, file appears, form resets', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    // Swallow the "Uploaded ... successfully" alert.
    const dialogPromise = captureDialog(page, 'accept');
    await opfs.uploadLocalFile(SMALL);
    await dialogPromise;

    await expect(opfs.row('small.txt')).toBeVisible();

    const stat = await statOPFSFile(page, 'small.txt');
    expect(stat.size).toBeGreaterThan(0);

    await expect(opfs.uploadName).toHaveValue('');
    await expect(opfs.uploadBtn).toBeDisabled();
  });

  test('C3: overwrite confirmation: accepting overwrites', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    // First upload (one info alert).
    const first = captureDialog(page, 'accept');
    await opfs.uploadLocalFile(SMALL);
    await first;

    // Second upload triggers confirm + info alert.
    const log = recordDialogs(page, 'accept');
    await opfs.uploadLocalFile(SMALL);
    // Wait until the second info alert has been observed.
    await expect.poll(() => log.messages.length).toBeGreaterThanOrEqual(2);
    log.stop();

    const text = await readOPFSFile(page, 'small.txt');
    expect(text).toContain('hello');
  });

  test('C4: overwrite confirmation: dismissing leaves the original intact', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    // Pre-seed via the editor so we know the expected content.
    await opfs.saveEditor('small.txt', 'ORIGINAL CONTENT');
    await expect(opfs.row('small.txt')).toBeVisible();

    // Now upload from disk and dismiss the overwrite confirm.
    const dialogPromise = captureDialog(page, 'dismiss');
    await opfs.uploadLocalFile(SMALL);
    await dialogPromise;

    const text = await readOPFSFile(page, 'small.txt');
    expect(text).toBe('ORIGINAL CONTENT');
  });

  test('C5: clearing the file input disables Upload again', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.uploadFileInput.setInputFiles(SMALL);
    await expect(opfs.uploadBtn).toBeEnabled();

    await opfs.uploadFileInput.setInputFiles([]);
    await expect(opfs.uploadName).toHaveValue('');
    await expect(opfs.uploadBtn).toBeDisabled();
  });

  test('C6: uploading under a different name uses that name', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    const dialogPromise = captureDialog(page, 'accept');
    await opfs.uploadLocalFile(SMALL, 'renamed.txt');
    await dialogPromise;

    await expect(opfs.row('renamed.txt')).toBeVisible();
    await expect(opfs.row('small.txt')).toHaveCount(0);
  });
});
