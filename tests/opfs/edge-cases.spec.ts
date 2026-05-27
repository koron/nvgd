import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:9280';

test.describe('OPFS Edge Cases', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/opfs/`);
  });

  test('empty directory name shows alert', async ({ page }) => {
    const [dialog] = await Promise.all([
      page.waitForEvent('dialog'),
      page.click('#mkdir-mkdir'),
    ]);
    expect(dialog.message()).toContain('Need directory name');
    await dialog.accept();
  });

  test('empty file name shows alert', async ({ page }) => {
    await page.fill('#editor-edit', 'content without name');
    const [dialog] = await Promise.all([
      page.waitForEvent('dialog'),
      page.click('#editor-save'),
    ]);
    expect(dialog.message()).toContain('Need file name');
    await dialog.accept();
  });

  test('overwrite prompt on duplicate upload', async ({ page }) => {
    await page.setInputFiles('#upload-file', {
      name: 'duplicate.txt',
      mimeType: 'text/plain',
      buffer: Buffer.from('first upload'),
    });
    await page.click('#upload-upload');
    await expect(page.locator('.grid-row')).toContainText('duplicate.txt');

    await page.setInputFiles('#upload-file', {
      name: 'duplicate.txt',
      mimeType: 'text/plain',
      buffer: Buffer.from('second upload'),
    });

    const [dialog] = await Promise.all([
      page.waitForEvent('dialog'),
      page.click('#upload-upload'),
    ]);
    expect(dialog.message()).toContain('overwrite');
    expect(dialog.message()).toContain('duplicate.txt');
    await dialog.accept();
  });

  test('overwrite prompt on duplicate download URL', async ({ page }) => {
    await page.fill('#download-url', `${BASE_URL}/version://`);
    await page.fill('#download-as', 'version-dup.txt');
    await page.click('#download-download');
    await expect(page.locator('.grid-row')).toContainText('version-dup.txt');

    await page.fill('#download-url', `${BASE_URL}/version://`);
    await page.fill('#download-as', 'version-dup.txt');

    const [dialog] = await Promise.all([
      page.waitForEvent('dialog'),
      page.click('#download-download'),
    ]);
    expect(dialog.message()).toContain('overwrite');
    expect(dialog.message()).toContain('version-dup.txt');
    await dialog.accept();
  });
});
