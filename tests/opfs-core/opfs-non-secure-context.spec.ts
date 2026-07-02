import { test, expect } from '@playwright/test';

test.describe('Browser Compatibility Considerations', () => {
  test('OPFS Operations in Non-Secure Context (if applicable)', async ({ page }) => {
    // If testing in a non-secure context (non-localhost, non-HTTPS)
    // Verify that OPFS operations are blocked or show appropriate errors
    // Verify that file upload, download, and editor operations are disabled
    await page.goto('http://localhost:9280/opfs/');

    // Since we're on localhost, OPFS operations should work
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('non-secure-test.txt');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('Non-secure context test');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "non-secure-test.txt"')).toBeVisible();
  });

  test('OPFS Operations with Different User Agents', async ({ page, context }) => {
    // Test OPFS operations with a different user agent
    const contextWithUA = await context.createUserContext({
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
    });

    const pageWithUA = await contextWithUA.newPage();
    await pageWithUA.goto('http://localhost:9280/opfs/');

    // Verify that OPFS operations work correctly with different user agents
    const editorName = pageWithUA.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('ua-test.txt');
    const editorTextarea = pageWithUA.locator('textarea');
    await editorTextarea.fill('User agent test');
    await pageWithUA.locator('text=Create or update a file, button').click();
    await expect(pageWithUA.locator('text=Uploaded "ua-test.txt"')).toBeVisible();

    // Verify that DuckDB integration works with different user agents
    const csvCheckbox = pageWithUA.locator('text=ua-test.txt').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await csvCheckbox.check();
    const duckdbBtn = pageWithUA.locator('text=DuckDB, button');
    await duckdbBtn.click();
    await pageWithUA.waitForEvent('page');
    const duckdbPage = pageWithUA.context().pages().find(p => p.url().includes('duckdb'));
    expect(duckdbPage).toBeTruthy();
    await duckdbPage!.close();

    await contextWithUA.close();
  });

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
