import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:9280';

test.describe('OPFS File I/O', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/opfs/`);
  });

  test('create a text file with editor', async ({ page }) => {
    await page.fill('#editor-name', 'hello.txt');
    await page.fill('#editor-edit', 'Hello, OPFS!');
    await page.click('#editor-save');

    await expect(page.locator('.grid-row')).toContainText('hello.txt');
    await expect(page.locator('#editor-name')).toHaveValue('');
    await expect(page.locator('#editor-edit')).toHaveValue('');
  });

  test('load file into editor on edit action', async ({ page }) => {
    await page.fill('#editor-name', 'editme.txt');
    await page.fill('#editor-edit', 'content to edit');
    await page.click('#editor-save');
    await expect(page.locator('.grid-row')).toContainText('editme.txt');

    await page.locator('.grid-row a', { hasText: 'Edit' }).click();

    await expect(page.locator('#editor-name')).toHaveValue('editme.txt');
    await expect(page.locator('#editor-edit')).toHaveValue('content to edit');
  });

  test('modify file content and save', async ({ page }) => {
    await page.fill('#editor-name', 'modify.txt');
    await page.fill('#editor-edit', 'original');
    await page.click('#editor-save');
    await expect(page.locator('.grid-row')).toContainText('modify.txt');

    await page.locator('.grid-row a', { hasText: 'Edit' }).click();
    await page.fill('#editor-edit', 'modified content');
    await page.click('#editor-save');

    await page.locator('.grid-row a', { hasText: 'Edit' }).click();
    await expect(page.locator('#editor-edit')).toHaveValue('modified content');
  });

  test('tab key in textarea inserts tab character', async ({ page }) => {
    await page.fill('#editor-name', 'tabtest.txt');
    await page.fill('#editor-edit', 'before');
    await page.locator('#editor-edit').press('Tab');
    await page.locator('#editor-edit').fill('before\t');
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
