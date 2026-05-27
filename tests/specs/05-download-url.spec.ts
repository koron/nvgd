import { expect, test } from '@playwright/test';
import { OpfsPage } from '../pages/OpfsPage';
import { captureDialog, recordDialogs } from '../helpers/dialogs';
import { clearOPFS, existsOPFS, readOPFSFile } from '../helpers/opfs';

/**
 * The /version/ endpoint of the running nvgd is a tiny text body and
 * is always available, so we use it as a download source that doesn't
 * depend on extra config.
 */
const VERSION_URL = '/version/';

test.describe('E. URL download', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/opfs/');
    await clearOPFS(page);
  });

  test('E1: Download/Clear buttons start disabled', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await expect(opfs.downloadBtn).toBeDisabled();
    await expect(opfs.downloadClearBtn).toBeDisabled();
  });

  test('E2: http(s) URL + name enables the Download button', async ({
    page,
    baseURL,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.downloadUrl.fill(`${baseURL}${VERSION_URL}`);
    await expect(opfs.downloadBtn).toBeDisabled();

    await opfs.downloadAs.fill('version.txt');
    await expect(opfs.downloadBtn).toBeEnabled();
    await expect(opfs.downloadClearBtn).toBeEnabled();
  });

  test('E3: non-http(s) protocols keep Download disabled', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.downloadUrl.fill('ftp://example.com/x.txt');
    await opfs.downloadAs.fill('x.txt');
    await expect(opfs.downloadBtn).toBeDisabled();
  });

  test('E4: successful download writes the file into OPFS', async ({
    page,
    baseURL,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.downloadFromUrl(`${baseURL}${VERSION_URL}`, 'version.txt');

    await expect(opfs.row('version.txt')).toBeVisible();
    expect(await existsOPFS(page, 'version.txt')).toBe('file');
    const body = await readOPFSFile(page, 'version.txt');
    expect(body.length).toBeGreaterThan(0);
  });

  test('E5: overwrite confirm fires when target name already exists', async ({
    page,
    baseURL,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.downloadFromUrl(`${baseURL}${VERSION_URL}`, 'dup.txt');
    await expect(opfs.row('dup.txt')).toBeVisible();

    // Second download with the same name → confirm dialog.
    const dialogPromise = captureDialog(page, 'accept');
    await opfs.downloadFromUrl(`${baseURL}${VERSION_URL}`, 'dup.txt');
    const dialog = await dialogPromise;
    expect(dialog.message()).toMatch(/overwrite/i);
  });

  test('E6: failing URL surfaces an alert and writes nothing', async ({
    page,
    baseURL,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    // /file:///__definitely_does_not_exist returns an error response.
    const log = recordDialogs(page, 'accept');
    await opfs.downloadFromUrl(
      `${baseURL}/file:///__definitely_does_not_exist__`,
      'oops.txt',
    );
    // Wait for any alert to appear or for the row to NOT appear.
    await page.waitForTimeout(500);
    log.stop();

    expect(await existsOPFS(page, 'oops.txt')).toBeNull();
  });

  test('E7: Clear resets URL and name inputs', async ({ page, baseURL }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.downloadUrl.fill(`${baseURL}${VERSION_URL}`);
    await opfs.downloadAs.fill('v.txt');
    await expect(opfs.downloadBtn).toBeEnabled();

    await opfs.downloadClearBtn.click();
    await expect(opfs.downloadUrl).toHaveValue('');
    await expect(opfs.downloadAs).toHaveValue('');
    await expect(opfs.downloadBtn).toBeDisabled();
    await expect(opfs.downloadClearBtn).toBeDisabled();
  });
});
