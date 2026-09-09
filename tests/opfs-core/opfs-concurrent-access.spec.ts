import { test, expect } from '@playwright/test';

test.describe('Edge Cases and Error Handling', () => {
  test('OPFS Operations with Concurrent Access', async ({ page, browser }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a file in OPFS
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('concurrent.txt');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('Concurrent access test');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "concurrent.txt"')).toBeVisible();

    // Step 2: Open the file in DuckDB shell
    const csvCheckbox = page.locator('text=concurrent.txt').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await csvCheckbox.check();
    const duckdbBtn = page.locator('text=DuckDB, button');
    await duckdbBtn.click();
    await page.waitForEvent('page');
    const duckdbPage = page.context().pages().find(p => p.url().includes('duckdb'));
    expect(duckdbPage).toBeTruthy();

    // Step 3: Try to edit the file using the OPFS editor
    await editorName.fill('concurrent.txt');
    await page.locator('text=Create or update a file, button').click();
    const lockAlert = page.locator('text=FILE MAY BE LOCKED, text=concurrent');
    await expect(lockAlert).toBeVisible();

    // Step 4: Close the DuckDB tab
    await duckdbPage!.close();

    // Step 5: After closing DuckDB, verify the file is editable again
    await editorName.fill('concurrent.txt');
    await page.locator('text=Create or update a file, button').click();
    // Should not show lock alert anymore
    await expect(lockAlert).not.toBeVisible();
  });
});
