import { test, expect } from '@playwright/test';

test.describe('OPFS File Listing', () => {
  test('File Handle Locked by DuckDB', async ({ page, browser }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a CSV file in OPFS
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('locked.csv');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('name,value\nA,1\nB,2');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "locked.csv"')).toBeVisible();

    // Step 2: Open the CSV file in DuckDB shell
    const csvCheckbox = page.locator('text=locked.csv').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await csvCheckbox.check();
    const duckdbBtn = page.locator('text=DuckDB, button');
    await duckdbBtn.click();
    await page.waitForEvent('page');
    const duckdbPage = page.context().pages().find(p => p.url().includes('duckdb'));
    expect(duckdbPage).toBeTruthy();

    // Step 3: Try to edit the file using the OPFS editor
    await editorName.fill('locked.csv');
    await page.locator('text=Create or update a file, button').click();
    const lockAlert = page.locator('text=FILE MAY BE LOCKED, text=locked');
    await expect(lockAlert).toBeVisible();

    // Step 4: Try to delete the file while it is open in DuckDB
    const fileCheckbox = page.locator('text=locked.csv').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await fileCheckbox.check();
    const deleteBtn = page.locator('text=Delete, button');
    await deleteBtn.click();
    const deleteError = page.locator('text=Error, text=locked');
    await expect(deleteError).toBeVisible();

    // Step 5: Close the DuckDB tab
    await duckdbPage!.close();

    // Step 6: After closing DuckDB, verify the file is editable again
    await editorName.fill('locked.csv');
    await page.locator('text=Create or update a file, button').click();
    // Should not show lock alert anymore
    await expect(lockAlert).not.toBeVisible();
  });
});
