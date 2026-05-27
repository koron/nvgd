import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:9280';

test.describe('OPFS Bulk Operations', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/opfs/`);
  });

  test('delete a single file', async ({ page }) => {
    await page.fill('#editor-name', 'todelete.txt');
    await page.fill('#editor-edit', 'delete me');
    await page.click('#editor-save');
    await expect(page.locator('.grid-row')).toContainText('todelete.txt');

    await page.locator('input.selectedFile').check();

    page.once('dialog', dialog => {
      expect(dialog.message()).toContain('delete');
      dialog.accept();
    });
    await page.click('#command-delete');

    await expect(page.locator('.grid-row')).not.toContainText('todelete.txt');
  });

  test('cancel deletion keeps the file', async ({ page }) => {
    await page.fill('#editor-name', 'keepfile.txt');
    await page.fill('#editor-edit', 'keep me');
    await page.click('#editor-save');
    await expect(page.locator('.grid-row')).toContainText('keepfile.txt');

    await page.locator('input.selectedFile').check();

    page.once('dialog', dialog => {
      dialog.dismiss();
    });
    await page.click('#command-delete');

    await expect(page.locator('.grid-row')).toContainText('keepfile.txt');
  });

  test('delete multiple selected files', async ({ page }) => {
    for (let i = 0; i < 3; i++) {
      await page.fill('#editor-name', `file${i}.txt`);
      await page.fill('#editor-edit', `content ${i}`);
      await page.click('#editor-save');
      await expect(page.locator('.grid-row')).toContainText(`file${i}.txt`);
    }

    await page.click('#toggle-selection-all');

    page.once('dialog', dialog => {
      dialog.accept();
    });
    await page.click('#command-delete');

    await expect(page.locator('.grid-row')).toHaveCount(0);
  });

  test('toggle-all selects and deselects all files', async ({ page }) => {
    for (let i = 0; i < 3; i++) {
      await page.fill('#editor-name', `file${i}.txt`);
      await page.fill('#editor-edit', `content ${i}`);
      await page.click('#editor-save');
    }

    await expect(page.locator('input.selectedFile:checked')).toHaveCount(0);

    await page.click('#toggle-selection-all');
    await expect(page.locator('input.selectedFile:checked')).toHaveCount(3);

    await page.click('#toggle-selection-all');
    await expect(page.locator('input.selectedFile:checked')).toHaveCount(0);
  });

  test('toggle-all shows indeterminate when some selected', async ({ page }) => {
    for (let i = 0; i < 3; i++) {
      await page.fill('#editor-name', `file${i}.txt`);
      await page.fill('#editor-edit', `content ${i}`);
      await page.click('#editor-save');
    }

    await page.locator('input.selectedFile').nth(0).check();

    const isIndeterminate = await page.locator('#toggle-selection-all')
      .evaluate(el => el.indeterminate);
    expect(isIndeterminate).toBe(true);
  });

  test('delete a directory with contents', async ({ page }) => {
    await page.fill('#mkdir-name', 'parent');
    await page.click('#mkdir-mkdir');
    await expect(page.locator('.grid-row')).toContainText('parent/');

    await page.locator('.grid-row a', { hasText: 'parent/' }).click();

    await page.fill('#editor-name', 'child.txt');
    await page.fill('#editor-edit', 'child content');
    await page.click('#editor-save');
    await expect(page.locator('.grid-row')).toContainText('child.txt');

    await page.locator('#header a').first().click();

    await page.locator('input.selectedFile').check();

    page.once('dialog', dialog => {
      dialog.accept();
    });
    await page.click('#command-delete');

    await expect(page.locator('.grid-row')).not.toContainText('parent');
  });

  test('duckdb button enabled when files selected', async ({ page }) => {
    await expect(page.locator('#command-duckdb')).toBeDisabled();

    await page.setInputFiles('#upload-file', {
      name: 'data.csv',
      mimeType: 'text/csv',
      buffer: Buffer.from('id,name\n1,test'),
    });
    await page.click('#upload-upload');
    await expect(page.locator('.grid-row')).toContainText('data.csv');

    await page.locator('input.selectedFile').check();
    await expect(page.locator('#command-duckdb')).toBeEnabled();
  });

  test('open with duckdb creates new page', async ({ page, context }) => {
    await page.setInputFiles('#upload-file', {
      name: 'query.csv',
      mimeType: 'text/csv',
      buffer: Buffer.from('x,y\n1,2'),
    });
    await page.click('#upload-upload');
    await expect(page.locator('.grid-row')).toContainText('query.csv');

    await page.locator('input.selectedFile').check();

    const newPagePromise = context.waitForEvent('page');
    await page.click('#command-duckdb');
    const newPage = await newPagePromise;
    await newPage.waitForLoadState();

    expect(newPage.url()).toContain('/duckdb/');
    expect(newPage.url()).toContain('opfs=');
  });
});
