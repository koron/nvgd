import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:9280';

test.describe('OPFS Directory Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/opfs/`);
  });

  test('create a directory', async ({ page }) => {
    await page.fill('#mkdir-name', 'mydir');
    await page.click('#mkdir-mkdir');

    await expect(page.locator('.grid-row')).toContainText('mydir/');
  });

  test('navigate into directory and back via breadcrumb', async ({ page }) => {
    await page.fill('#mkdir-name', 'subdir');
    await page.click('#mkdir-mkdir');
    await expect(page.locator('.grid-row')).toContainText('subdir/');

    await page.locator('.grid-row a', { hasText: 'subdir/' }).click();
    await expect(page.locator('#header')).toContainText('subdir');

    await page.locator('#header a').first().click();
    await expect(page.locator('.grid-row')).toContainText('subdir/');
  });

  test('reload button refreshes listing', async ({ page }) => {
    await page.fill('#mkdir-name', 'reloadtest');
    await page.click('#mkdir-mkdir');
    await expect(page.locator('.grid-row')).toContainText('reloadtest/');

    await page.click('#command-reload');
    await expect(page.locator('.grid-row')).toContainText('reloadtest/');
  });

  test('deep directory nesting', async ({ page }) => {
    const dirs = ['a', 'b', 'c', 'd', 'e'];
    for (const d of dirs) {
      await page.fill('#mkdir-name', d);
      await page.click('#mkdir-mkdir');
      await expect(page.locator('.grid-row')).toContainText(`${d}/`);
      await page.locator('.grid-row a', { hasText: `${d}/` }).click();
    }

    await expect(page.locator('#header')).toContainText('e');
  });
});
