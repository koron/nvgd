import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:9280';

test.describe('OPFS Rendering', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/opfs/`);
  });

  test('page title shows OPFS', async ({ page }) => {
    await expect(page).toHaveTitle(/OPFS/);
  });

  test('all UI sections are rendered', async ({ page }) => {
    await expect(page.locator('#header')).toBeVisible();
    await expect(page.locator('#command-reload')).toBeVisible();
    await expect(page.locator('#command-delete')).toBeVisible();
    await expect(page.locator('#command-duckdb')).toBeVisible();
    await expect(page.locator('#mkdir-name')).toBeVisible();
    await expect(page.locator('#mkdir-mkdir')).toBeVisible();
    await expect(page.locator('#upload-file')).toBeVisible();
    await expect(page.locator('#upload-name')).toBeVisible();
    await expect(page.locator('#upload-upload')).toBeVisible();
    await expect(page.locator('#editor-name')).toBeVisible();
    await expect(page.locator('#editor-edit')).toBeVisible();
    await expect(page.locator('#editor-save')).toBeVisible();
    await expect(page.locator('#editor-clear')).toBeVisible();
    await expect(page.locator('#download-url')).toBeVisible();
    await expect(page.locator('#download-as')).toBeVisible();
    await expect(page.locator('#download-download')).toBeVisible();
  });

  test('delete and duckdb buttons are initially disabled', async ({ page }) => {
    await expect(page.locator('#command-delete')).toBeDisabled();
    await expect(page.locator('#command-duckdb')).toBeDisabled();
  });
});
