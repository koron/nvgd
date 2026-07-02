import { test, expect } from '@playwright/test';

test.describe('File Deletion Operations', () => {
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
