import { test, expect } from '@playwright/test';

test.describe('Browser Compatibility Considerations', () => {
  test('OPFS Operations with Different Screen Sizes', async ({ page }) => {
    // Test OPFS operations at different screen sizes
    // Desktop size
    await page.setViewportSize({ width: 1920, height: 1080 });
    await page.goto('http://localhost:9280/opfs/');
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('desktop-test.txt');
    await page.locator('textarea').fill('Desktop test');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "desktop-test.txt"')).toBeVisible();

    // Tablet size
    await page.setViewportSize({ width: 768, height: 1024 });
    await expect(page.locator('text=desktop-test.txt')).toBeVisible();

    // Mobile size
    await page.setViewportSize({ width: 375, height: 667 });
    await expect(page.locator('text=desktop-test.txt')).toBeVisible();

    // Verify all OPFS operations work at different screen sizes
    await expect(page.locator('text=Reload, button, :text("Reload")')).toBeVisible();
    await expect(page.locator('text=Delete, button')).toBeVisible();
    await expect(page.locator('text=DuckDB, button')).toBeVisible();
  });
});
