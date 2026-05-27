import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:9280';

test.describe('OPFS Edge Cases', () => {
  test.skip(({ browserName }) => browserName === 'webkit', 'OPFS API requires secure context not available in headless WebKit');

  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/opfs/`);
  });

  test('empty directory name shows alert', async ({ page }) => {
    let dialogMessage = '';
    page.on('dialog', dialog => {
      dialogMessage = dialog.message();
      dialog.accept();
    });
    await page.click('#mkdir-mkdir');
    expect(dialogMessage).toContain('Need directory name');
  });

  test('empty file name shows alert', async ({ page }) => {
    await page.fill('#editor-edit', 'content without name');
    let dialogMessage = '';
    page.on('dialog', dialog => {
      dialogMessage = dialog.message();
      dialog.accept();
    });
    await page.click('#editor-save');
    expect(dialogMessage).toContain('Need file name');
  });

  test('overwrite prompt on duplicate upload', async ({ page }) => {
    let confirmMessage = '';
    page.on('dialog', dialog => {
      if (dialog.message().includes('overwrite')) {
        confirmMessage = dialog.message();
      }
      dialog.accept();
    });

    await page.setInputFiles('#upload-file', {
      name: 'duplicate.txt',
      mimeType: 'text/plain',
      buffer: Buffer.from('first upload'),
    });
    await page.click('#upload-upload');
    await expect(page.locator('.grid-row').filter({ hasText: 'duplicate.txt' })).toBeVisible();

    await page.setInputFiles('#upload-file', {
      name: 'duplicate.txt',
      mimeType: 'text/plain',
      buffer: Buffer.from('second upload'),
    });

    await page.click('#upload-upload');
    await expect.poll(() => confirmMessage).toContain('overwrite');
    expect(confirmMessage).toContain('duplicate.txt');
  });

  test('overwrite prompt on duplicate download URL', async ({ page }) => {
    let confirmMessage = '';
    page.on('dialog', dialog => {
      if (dialog.message().includes('overwrite')) {
        confirmMessage = dialog.message();
      }
      dialog.accept();
    });

    await page.fill('#download-url', `${BASE_URL}/version://`);
    await page.fill('#download-as', 'version-dup.txt');
    await page.click('#download-download');
    await expect(page.locator('.grid-row').filter({ hasText: 'version-dup.txt' })).toBeVisible();

    await page.fill('#download-url', `${BASE_URL}/version://`);
    await page.fill('#download-as', 'version-dup.txt');

    await page.click('#download-download');
    await expect.poll(() => confirmMessage).toContain('overwrite');
    expect(confirmMessage).toContain('version-dup.txt');
  });
});
