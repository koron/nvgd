import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:9280';

test.describe('OPFS Download URL', () => {
  test.skip(({ browserName }) => browserName === 'webkit', 'OPFS API requires secure context not available in headless WebKit');
  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/opfs/`);
  });

  test('download button enabled only with valid inputs', async ({ page }) => {
    await expect(page.locator('#download-download')).toBeDisabled();

    await page.fill('#download-url', 'http://localhost:9280/version://');
    await expect(page.locator('#download-download')).toBeDisabled();

    await page.fill('#download-as', 'version.txt');
    await expect(page.locator('#download-download')).toBeEnabled();
  });

  test('clear button resets download inputs', async ({ page }) => {
    await expect(page.locator('#download-clear')).toBeDisabled();
    await expect(page.locator('#download-clear')).toBeDisabled();

    await page.fill('#download-url', 'http://localhost:9280/');
    await expect(page.locator('#download-clear')).toBeEnabled();

    await page.click('#download-clear');
    await expect(page.locator('#download-url')).toHaveValue('');
    await expect(page.locator('#download-as')).toHaveValue('');
  });

  test('download from NVGD URL saves file to OPFS', async ({ page }) => {
    await page.fill('#download-url', 'http://localhost:9280/version://');
    await page.fill('#download-as', 'downloaded-version.txt');
    await page.click('#download-download');

    await expect(page.locator('.grid-row')).toContainText('downloaded-version.txt');
  });
});
