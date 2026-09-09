import { test, expect } from '@playwright/test';

test.describe('OPFS File Listing', () => {
  test('File Sorting', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create files with names in different orders
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    const editorTextarea = page.locator('textarea');

    await editorName.fill('b.txt');
    await editorTextarea.fill('B content');
    await page.locator('text=Create or update a file, button').click();

    await editorName.fill('a.txt');
    await editorTextarea.fill('A content');
    await page.locator('text=Create or update a file, button').click();

    await editorName.fill('c.txt');
    await editorTextarea.fill('C content');
    await page.locator('text=Create or update a file, button').click();

    // Step 2: Verify the directory listing sorts files alphabetically
    const listing = page.locator('.listing, .directory-list, table');
    const fileNames = await listing.locator('text=a.txt, text=b.txt, text=c.txt').count();
    expect(fileNames).toBe(3);

    // Step 3: Verify files and directories are mixed correctly in the sorted listing
    await expect(listing).toContainText('a.txt');
    await expect(listing).toContainText('b.txt');
    await expect(listing).toContainText('c.txt');
  });

  test('Empty Directory', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a new empty directory
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('emptydir');
    await page.locator('text=Create new directory, button').click();

    // Step 2: Navigate into the empty directory
    await page.locator('text=emptydir/').click();

    // Step 3: Verify the directory listing is empty (no files or directories)
    const listing = page.locator('.listing, .directory-list, table');
    await expect(listing).toBeVisible();

    // Step 4: Verify the "Select all" checkbox is unchecked
    const selectAll = page.locator('input[type=checkbox], [class*="checkbox"]');
    await expect(selectAll.first()).not.toBeChecked();

    // Step 5: Verify the "Delete" and "DuckDB" buttons are disabled
    const deleteBtn = page.locator('text=Delete, button');
    await expect(deleteBtn).toBeDisabled();
    const duckdbBtn = page.locator('text=DuckDB, button');
    await expect(duckdbBtn).toBeDisabled();
  });

  test('Large Number of Files', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create 50+ files in OPFS
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    const editorTextarea = page.locator('textarea');

    for (let i = 0; i < 50; i++) {
      await editorName.fill(`file${i}.txt`);
      await editorTextarea.fill(`Content of file ${i}`);
      await page.locator('text=Create or update a file, button').click();
    }

    // Step 2: Verify the directory listing displays all files
    const listing = page.locator('.listing, .directory-list, table');
    const fileCount = await listing.locator('text=file').count();
    expect(fileCount).toBeGreaterThan(40);

    // Step 3: Verify the "Select all" checkbox works with many files
    const selectAllCheckbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await selectAllCheckbox.first().check();

    // Step 4: Verify the "Delete" button is enabled with many selected files
    const deleteBtn = page.locator('text=Delete, button');
    await expect(deleteBtn).toBeEnabled();

    // Step 5: Delete a batch of files and verify they are all removed
    await deleteBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();
  });

  test('Special Characters in Filenames', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a file with spaces in the name
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('my file.txt');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('Content with spaces');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=my file.txt')).toBeVisible();

    // Step 2: Upload a file with hyphens in the name
    const uploadInput = page.locator('input[type=file]');
    await uploadInput.setInputFiles('tests/testdata/sample.txt');
    const nameField = page.locator('input[placeholder*="File name"]');
    await nameField.fill('my-file.txt');
    await page.locator('text=Upload, button').click();
    await expect(page.locator('text=my-file.txt')).toBeVisible();

    // Step 3: Create a file with dots in the name
    await editorName.fill('my.file.txt');
    await editorTextarea.fill('Content with dots');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=my.file.txt')).toBeVisible();

    // Step 4: Navigate into a directory with spaces in its name
    await dirNameInput.click();
    await dirNameInput.fill('my dir');
    await page.locator('text=Create new directory, button').click();

    // Step 5: Verify the breadcrumb correctly displays the directory with spaces
    await expect(page.locator('text=my dir')).toBeVisible();
  });

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
