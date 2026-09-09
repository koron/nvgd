import { test, expect } from '@playwright/test';

test.describe('File Upload Operations', () => {
  test('Upload Local File', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Select a local file using the file picker
    const uploadInput = page.locator('input[type=file]');
    await uploadInput.setInputFiles('tests/testdata/sample.txt');

    // Step 2: Verify the file name populates the "File name to upload" field
    const nameField = page.locator('input[placeholder*="File name"]');
    await expect(nameField).toHaveValue(/sample\.txt/);

    // Step 3: Click "Upload" button
    await page.locator('text=Upload, button').click();

    // Step 4: Verify success alert
    await expect(page.locator('text=Uploaded "sample.txt"')).toBeVisible();

    // Step 5: Verify the uploaded file appears in the directory listing with correct size
    const listing = page.locator('.listing, .directory-list, table');
    await expect(listing).toContainText('sample.txt');

    // Step 6: Verify the file shows the correct last modified date
    await expect(page.locator('text=sample.txt')).toBeVisible();

    // Step 7: Upload an empty file and verify size is 0
    await uploadInput.setInputFiles('tests/testdata/empty.txt');
    await page.locator('text=Upload, button').click();
    await expect(page.locator('text=Uploaded "empty.txt"')).toBeVisible();
    await expect(page.locator('text=empty.txt')).toBeVisible();

    // Step 8: Upload a file with a name that already exists
    await uploadInput.setInputFiles('tests/testdata/sample.txt');
    await page.locator('text=Upload, button').click();

    // Step 9: Verify a confirmation dialog appears before overwriting
    const confirmDialog = page.locator('text=Confirm, text=overwrite');
    await expect(confirmDialog).toBeVisible();

    // Step 10: After confirming overwrite, verify the file listing shows updated size/date
    await page.locator('text=Confirm, button, :text("Confirm")').click();
    await expect(page.locator('text=sample.txt')).toBeVisible();

    // Step 11: Upload a file with an empty name
    await uploadInput.setInputFiles('tests/testdata/sample.txt');
    const nameField2 = page.locator('input[placeholder*="File name"]');
    await nameField2.fill('');
    await page.locator('text=Upload, button').click();

    // Step 12: Verify an alert `Need file name` is shown
    const alert = page.locator('text=Need file name');
    await expect(alert).toBeVisible();
  });

  test('Upload File with Custom Name', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Select a local file `source.txt`
    const uploadInput = page.locator('input[type=file]');
    await uploadInput.setInputFiles('tests/testdata/sample.txt');

    // Step 2: Change the file name in the "File name to upload" field to `dest.txt`
    const nameField = page.locator('input[placeholder*="File name"]');
    await nameField.fill('dest.txt');

    // Step 3: Click "Upload" button
    await page.locator('text=Upload, button').click();

    // Step 4: Verify the file appears under its new name `dest.txt`
    await expect(page.locator('text=dest.txt')).toBeVisible();

    // Step 5: Verify the original file `source.txt` does not appear in the listing
    await expect(page.locator('text=source.txt')).not.toBeVisible();
  });
});
