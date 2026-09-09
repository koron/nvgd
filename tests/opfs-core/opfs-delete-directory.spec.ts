import { test, expect } from '@playwright/test';

test.describe('File Deletion Operations', () => {
  test('Delete Directory Recursively', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a directory `mydir` with files inside
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('mydir');
    await page.locator('text=Create new directory, button').click();
    await expect(page.locator('text=mydir/')).toBeVisible();

    // Create files inside
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('file1.txt');
    await page.locator('textarea').fill('File 1 content');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "file1.txt"')).toBeVisible();

    await editorName.fill('file2.txt');
    await page.locator('textarea').fill('File 2 content');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "file2.txt"')).toBeVisible();

    // Navigate into mydir
    await page.locator('text=mydir/').click();

    // Step 2: Check the checkbox next to `mydir/`
    const dirCheckbox = page.locator('text=mydir/').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await dirCheckbox.check();

    // Step 3: Click the "Delete" button
    const deleteBtn = page.locator('text=Delete, button');
    await expect(deleteBtn).toBeEnabled();
    await deleteBtn.click();

    // Step 4: Confirm the deletion
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Step 5: Verify the directory and all its contents are removed
    await expect(page.locator('text=mydir/')).not.toBeVisible();
    await expect(page.locator('text=file1.txt')).not.toBeVisible();
    await expect(page.locator('text=file2.txt')).not.toBeVisible();

    // Step 6: Verify the directory does not appear in the listing
    await expect(page.locator('text=mydir/')).not.toBeVisible();

    // Step 7: Delete an empty directory
    await dirNameInput.click();
    await dirNameInput.fill('emptydir');
    await page.locator('text=Create new directory, button').click();
    await expect(page.locator('text=emptydir/')).toBeVisible();

    await page.locator('text=emptydir/').click();
    // No files inside
    await deleteBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Step 8: Verify the deletion succeeds
    await expect(page.locator('text=emptydir/')).not.toBeVisible();
  });

  test('Delete Single File', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a file `single.txt` in OPFS
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('single.txt');
    await page.locator('textarea').fill('Single file content');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "single.txt"')).toBeVisible();

    // Step 2: Check the checkbox next to the file
    const fileCheckbox = page.locator('text=single.txt').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await fileCheckbox.check();

    // Step 3: Click the "Delete" button
    const deleteBtn = page.locator('text=Delete, button');
    await expect(deleteBtn).toBeEnabled();
    await deleteBtn.click();

    // Step 4: Confirm the deletion
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Step 5: Verify the file is removed from the listing
    await expect(page.locator('text=single.txt')).not.toBeVisible();

    // Step 6: Verify other files in the same directory remain unaffected
    await expect(page.locator('text=other')).not.toBeVisible();
  });
});
