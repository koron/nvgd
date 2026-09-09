import { test, expect } from '@playwright/test';

test.describe('DuckDB Integration', () => {
  test('Unsupported File Types', async ({ page, browser }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a binary file (e.g., a small dummy binary) in OPFS
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('binary.dat');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('BINARY_DATA_CONTENT_HERE');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "binary.dat"')).toBeVisible();

    // Step 2: Check the checkbox next to the unsupported file
    const binaryCheckbox = page.locator('text=binary.dat').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await binaryCheckbox.check();

    // Step 3: Click the "DuckDB" button
    const duckdbBtn = page.locator('text=DuckDB, button');
    await duckdbBtn.click();

    // Step 4: Verify the DuckDB shell opens but does not create a view for the unsupported file
    await page.waitForEvent('page');
    const duckdbPage = page.context().pages().find(p => p.url().includes('duckdb'));
    expect(duckdbPage).toBeTruthy();

    // The view should not be named opfs0 for this unsupported file
    const viewOpfs0 = duckdbPage!.locator('text=opfs0');
    await expect(viewOpfs0).not.toBeVisible();

    // Step 5: Verify only supported file types generate DuckDB views
    await expect(duckdbPage!.locator('text=binary.dat')).toBeVisible();

    // Close DuckDB tab
    await duckdbPage!.close();
  });
});
