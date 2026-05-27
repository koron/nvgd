import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:9280';

test.describe('OPFS Bulk Operations', () => {
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

  test('delete a single file', async ({ page }) => {
    await createFile(page, 'todelete.txt', 'delete me');

    await page.locator('input.selectedFile').check();

    page.on('dialog', dialog => {
      expect(dialog.message()).toContain('delete');
      dialog.accept();
    });
    await page.click('#command-delete');

    await expect(page.locator('.grid-row').filter({ hasText: 'todelete.txt' })).toHaveCount(0);
  });

  test('cancel deletion keeps the file', async ({ page }) => {
    await createFile(page, 'keepfile.txt', 'keep me');

    await page.locator('input.selectedFile').check();

    page.on('dialog', dialog => {
      dialog.dismiss();
    });
    await page.click('#command-delete');

    await expect(page.locator('.grid-row').filter({ hasText: 'keepfile.txt' })).toBeVisible();
  });

  test('delete multiple selected files', async ({ page }) => {
    for (let i = 0; i < 3; i++) {
      await createFile(page, `file${i}.txt`, `content ${i}`);
    }

    await page.click('#toggle-selection-all');

    page.on('dialog', dialog => {
      dialog.accept();
    });
    await page.click('#command-delete');

    await expect(page.locator('.grid-row')).toHaveCount(0);
  });

  test('toggle-all selects and deselects all files', async ({ page }) => {
    for (let i = 0; i < 3; i++) {
      await createFile(page, `file${i}.txt`, `content ${i}`);
    }

    await expect(page.locator('input.selectedFile:checked')).toHaveCount(0);

    await page.click('#toggle-selection-all');
    await expect(page.locator('input.selectedFile:checked')).toHaveCount(3);

    await page.click('#toggle-selection-all');
    await expect(page.locator('input.selectedFile:checked')).toHaveCount(0);
  });

  test('toggle-all shows indeterminate when some selected', async ({ page }) => {
    for (let i = 0; i < 3; i++) {
      await createFile(page, `file${i}.txt`, `content ${i}`);
    }

    await page.locator('input.selectedFile').nth(0).check();

    const isIndeterminate = await page.locator('#toggle-selection-all')
      .evaluate(el => el.indeterminate);
    expect(isIndeterminate).toBe(true);
  });

  test('delete a directory with contents', async ({ page }) => {
    await page.fill('#mkdir-name', 'parent');
    await page.click('#mkdir-mkdir');
    await expect(page.locator('.grid-row').filter({ hasText: 'parent/' })).toBeVisible();

    await page.locator('.grid-row a', { hasText: 'parent/' }).click();

    await createFile(page, 'child.txt', 'child content');

    await page.locator('#header a').first().click();

    await page.locator('input.selectedFile').check();

    page.on('dialog', dialog => {
      dialog.accept();
    });
    await page.click('#command-delete');

    await expect(page.locator('.grid-row').filter({ hasText: 'parent/' })).toHaveCount(0);
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
