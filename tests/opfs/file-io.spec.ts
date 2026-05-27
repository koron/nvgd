import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:9280';

test.describe('OPFS File I/O', () => {
  test.skip(({ browserName }) => browserName === 'webkit', 'OPFS API requires secure context not available in headless WebKit');

  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/opfs/`);
  });

  async function createFile(page, name, content) {
    await page.fill('#editor-name', name);
    await page.fill('#editor-edit', content);
    await page.click('#editor-save');
    await expect(page.locator('#editor-name')).toHaveValue('');
    await expect(page.locator('.grid-row').filter({ hasText: name })).toBeVisible();
  }

  test('create a text file with editor', async ({ page }) => {
    await createFile(page, 'hello.txt', 'Hello, OPFS!');
  });

  test('load file into editor on edit action', async ({ page }) => {
    await createFile(page, 'editme.txt', 'content to edit');

    await page.locator('.grid-row a', { hasText: 'Edit' }).click();

    await expect(page.locator('#editor-name')).toHaveValue('editme.txt');
    await expect(page.locator('#editor-edit')).toHaveValue('content to edit');
  });

  test('modify file content and save', async ({ page }) => {
    await createFile(page, 'modify.txt', 'original');

    await page.locator('.grid-row a', { hasText: 'Edit' }).click();
    await expect(page.locator('#editor-edit')).toHaveValue('original');
    await page.fill('#editor-edit', 'modified content');
    await page.click('#editor-save');
    await expect(page.locator('#editor-name')).toHaveValue('');

    await page.locator('.grid-row a', { hasText: 'Edit' }).click();
    await expect(page.locator('#editor-edit')).toHaveValue('modified content');
  });

  test('tab key in textarea inserts tab character', async ({ page }) => {
    await createFile(page, 'tabtest.txt', 'before');

    await page.locator('.grid-row a', { hasText: 'Edit' }).click();
    await expect(page.locator('#editor-edit')).toHaveValue('before');
    await page.locator('#editor-edit').press('Tab');
    await expect(page.locator('#editor-edit')).toHaveValue(/before\t/);
  });

  test('editor clear button resets fields', async ({ page }) => {
    await page.fill('#editor-name', 'somefile.txt');
    await page.fill('#editor-edit', 'some content');
    await page.click('#editor-clear');

    await expect(page.locator('#editor-name')).toHaveValue('');
    await expect(page.locator('#editor-edit')).toHaveValue('');
  });

  test('selecting file auto-fills name and enables upload button', async ({ page }) => {
    await page.setInputFiles('#upload-file', {
      name: 'auto-upload.txt',
      mimeType: 'text/plain',
      buffer: Buffer.from('auto upload test'),
    });

    await expect(page.locator('#upload-name')).toHaveValue('auto-upload.txt');
    await expect(page.locator('#upload-upload')).toBeEnabled();
  });

  test('upload a file appears in listing', async ({ page }) => {
    await page.setInputFiles('#upload-file', {
      name: 'uploaded.txt',
      mimeType: 'text/plain',
      buffer: Buffer.from('uploaded content'),
    });
    await page.click('#upload-upload');

    await expect(page.locator('.grid-row')).toContainText('uploaded.txt');
  });

  test('upload with custom name', async ({ page }) => {
    await page.setInputFiles('#upload-file', {
      name: 'original-name.txt',
      mimeType: 'text/plain',
      buffer: Buffer.from('custom name test'),
    });
    await page.fill('#upload-name', 'custom-name.txt');
    await page.click('#upload-upload');

    await expect(page.locator('.grid-row')).toContainText('custom-name.txt');
  });
});
