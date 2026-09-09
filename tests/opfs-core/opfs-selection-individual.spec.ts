import { test, expect } from '@playwright/test';

test.describe('File Selection Operations', () => {
  test('Select Individual Files', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create two files in OPFS
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    const editorTextarea = page.locator('textarea');

    await editorName.fill('alpha.txt');
    await editorTextarea.fill('Alpha content');
    await page.locator('text=Create or update a file, button').click();

    await editorName.fill('beta.txt');
    await editorTextarea.fill('Beta content');
    await page.locator('text=Create or update a file, button').click();

    // Step 2: Check the checkbox for the first file only
    const alphaCheckbox = page.locator('text=alpha.txt').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await alphaCheckbox.check();

    // Step 3: Verify the "Select all" checkbox is unchecked
    const selectAll = page.locator('input[type=checkbox], [class*="checkbox"]');
    await expect(selectAll.first()).not.toBeChecked();

    // Step 4: Check the checkbox for the second file
    const betaCheckbox = page.locator('text=beta.txt').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await betaCheckbox.check();

    // Step 5: Verify the "Select all" checkbox becomes indeterminate
    const selectAllState = await selectAll.first().isChecked();
    expect(selectAllState).toBeUndefined();

    // Step 6: Uncheck the second file
    await betaCheckbox.uncheck();

    // Step 7: Verify the "Select all" checkbox becomes unchecked
    await expect(selectAll.first()).not.toBeChecked();
  });

  test('Mixed Files and Directories', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a directory and a file in OPFS
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('mixeddir');
    await page.locator('text=Create new directory, button').click();

    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('mixed.txt');
    await page.locator('textarea').fill('Mixed content');
    await page.locator('text=Create or update a file, button').click();

    // Step 2: Check the checkbox for the file
    const fileCheckbox = page.locator('text=mixed.txt').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await fileCheckbox.check();

    // Step 3: Check the checkbox for the directory
    const dirCheckbox = page.locator('text=mixeddir/').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await dirCheckbox.check();

    // Step 4: Verify both are selected simultaneously
    await expect(page.locator('text=mixed.txt')).toBeVisible();
    await expect(page.locator('text=mixeddir/')).toBeVisible();

    // Step 5: Verify the "Delete" button is enabled
    const deleteBtn = page.locator('text=Delete, button');
    await expect(deleteBtn).toBeEnabled();

    // Step 6: Verify the "DuckDB" button is enabled
    const duckdbBtn = page.locator('text=DuckDB, button');
    await expect(duckdbBtn).toBeEnabled();
  });
});
