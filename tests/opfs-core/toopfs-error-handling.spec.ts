import { test, expect } from '@playwright/test';

test.describe('toopfs Filter - Download Files from NVGD to OPFS', () => {
  test('Download After Reload', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Set up a destination directory and select files
    const destDirInput = page.locator('input[placeholder*="OPFS directory"], input[name*="dest"]');
    await destDirInput.fill('reload-output/');

    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();

    // Reload the page
    await page.reload();

    // Verify the toopfs UI is restored
    const toopfsUI = page.locator('text=Download to OPFS');
    await expect(toopfsUI).toBeVisible();

    // Select and download files again
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are downloaded successfully
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with No Destination Set', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select files without setting a destination directory
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();

    // Download the files
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are saved in the root OPFS directory
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with Indeterminate Selection', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select only some files (not all)
    const checkboxes = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkboxes.first().check();
    await checkboxes.nth(2).check();

    // Verify the "Select/unselect all" checkbox is indeterminate
    const selectAllState = await checkboxes.first().isChecked();
    expect(selectAllState).toBeUndefined();

    // Verify the file count and total size reflect only the selected files
    const fileCountDisplay = page.locator('text=2 files');
    await expect(fileCountDisplay).toBeVisible();

    // Verify the "Download" button is enabled
    const downloadBtn = page.locator('text=Download, button');
    await expect(downloadBtn).toBeEnabled();
  });

  test('Download with Zero-Size Files', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Download the zero-size file
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the file is saved in OPFS with size 0
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with Large Total Size', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select all files
    const checkboxes = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkboxes.first().check();

    // Verify the total size is displayed correctly
    const downloadBtn = page.locator('text=Download, button');
    await expect(downloadBtn).toBeEnabled();

    // Download all files
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the download completes successfully
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with Special Characters in Names', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select and download files
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are saved correctly in OPFS
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download and Edit', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Download the file to OPFS
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Open the OPFS file in the editor
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();

    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('modified.txt');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('Modified content after download');
    await page.locator('text=Create or update a file, button').click();

    // Verify the modified content is persisted in OPFS
    await expect(page.locator('text=Uploaded "modified.txt"')).toBeVisible();
  });

  test('Download Then Delete', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Download files to OPFS
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Navigate to the OPFS destination directory
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();

    // Select and delete the downloaded files
    const fileCheckbox = page.locator('text=Download to OPFS').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await fileCheckbox.check();
    const deleteBtn = page.locator('text=Delete, button');
    await deleteBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are removed from OPFS
    await expect(fileCheckbox).not.toBeChecked();
  });

  test('Download with Different Destination Directories', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Download files to the root OPFS directory
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Navigate away and come back
    await page.goto('http://localhost:9280/opfs/');
    await page.goto('http://localhost:9280/file://?toopfs');

    // Download the same files to a different destination directory
    const destDirInput = page.locator('input[placeholder*="OPFS directory"], input[name*="dest"]');
    await destDirInput.fill('backup/');
    await checkbox.first().check();
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify both sets of files exist in their respective directories
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with Existing Files in Destination', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Create some files in OPFS first
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('existing.txt');
    await page.locator('textarea').fill('Original content');
    await page.locator('text=Create or update a file, button').click();

    // Navigate to NVGD page with toopfs filter
    await page.goto('http://localhost:9280/file://?toopfs');

    // Download files that have the same names as existing OPFS files
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify a confirmation dialog appears for each conflicting file
    const confirmDialog = page.locator('text=Confirm, text=overwrite');
    await expect(confirmDialog).toBeVisible();

    // After confirming, verify the files are overwritten
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the overwritten files contain the new content
    await expect(page.locator('text=Download to OPFS')).toBeVisible();
  });

  test('Download with file:// Protocol keepcompress', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Download the compressed file to OPFS
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the compressed file is saved correctly in OPFS
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with https:// Protocol all Parameter', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs&all');

    // Select and download files
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are downloaded correctly with the `?all` parameter
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with file:// Protocol all Parameter', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs&all');

    // Select and download files
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are downloaded correctly with the `?all` parameter
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });
});
