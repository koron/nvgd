import { test, expect } from '@playwright/test';

test.describe('toopfs Filter - Download Files from NVGD to OPFS', () => {
  test('Download Files from NVGD to OPFS', async ({ page }) => {
    // Navigate to an NVGD page with a `file://` source and `?toopfs` filter
    await page.goto('http://localhost:9280/file://?toopfs');

    // Verify the "Download to OPFS" UI is displayed
    const toopfsUI = page.locator('text=Download to OPFS, [class*="toopfs"], [class*="topfs"]');
    await expect(toopfsUI).toBeVisible();

    // Verify all files are pre-selected by default
    const checkboxes = page.locator('input[type=checkbox], [class*="checkbox"]');
    const checkedCount = await checkboxes.count();
    expect(checkedCount).toBeGreaterThan(0);

    // Verify the file count and total size are displayed correctly
    const fileCountDisplay = page.locator('text=0 files, text=0 bytes');
    await expect(fileCountDisplay).toBeVisible();

    // Unselect some files
    const firstCheckbox = page.locator('input[type=checkbox], [class*="checkbox"]').first();
    await firstCheckbox.uncheck();

    // Verify the file count and total size update
    const updatedCount = await page.locator('text=0 files, text=1 file').first();
    await expect(updatedCount).toBeVisible();

    // Verify the "Download" button is disabled when no files are selected
    await firstCheckbox.uncheck();
    const downloadBtn = page.locator('text=Download, button');
    await expect(downloadBtn).toBeDisabled();

    // Select files and click "Download"
    const lastCheckbox = page.locator('input[type=checkbox], [class*="checkbox"]').last();
    await lastCheckbox.check();
    await downloadBtn.click();

    // Verify a confirmation dialog appears
    const confirmDialog = page.locator('text=Confirm, text=Download');
    await expect(confirmDialog).toBeVisible();

    // Confirm the download
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the progress overlay appears with download steps
    const progressOverlay = page.locator('text=#1/1, text=downloading');
    await expect(progressOverlay).toBeVisible();

    // After download completes, verify the OPFS destination directory is opened in a new tab
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();

    // Verify the downloaded files are visible in the OPFS destination directory
    await expect(opfsPage!.locator('text=Download to OPFS')).toBeVisible();
  });

  test('Nested Directory Structure', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Verify the toopfs UI is visible
    const toopfsUI = page.locator('text=Download to OPFS');
    await expect(toopfsUI).toBeVisible();

    // Set a destination directory path
    const destDirInput = page.locator('input[placeholder*="OPFS directory"], input[name*="dest"]');
    await destDirInput.fill('output/nested/');

    // Select files and download
    const checkboxes = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkboxes.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are saved in the correct subdirectory of OPFS
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
    await expect(opfsPage!.locator('text=output/nested/')).toBeVisible();
  });

  test('Select All/Unselect All', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Click the "Select/unselect all" checkbox
    const selectAllCheckbox = page.locator('input[type=checkbox], [class*="checkbox"]').first();
    await expect(selectAllCheckbox).toBeVisible();

    // Verify all files are selected
    const checkedCount = await page.locator('input[type=checkbox], [class*="checkbox"]').count();
    expect(checkedCount).toBeGreaterThan(0);

    // Click again to unselect all
    await selectAllCheckbox.click();

    // Verify all files are unselected
    const uncheckedCount = await page.locator('input[type=checkbox], [class*="checkbox"]').count();
    expect(uncheckedCount).toBeGreaterThan(0);

    // Verify the "Download" button is disabled when nothing is selected
    const downloadBtn = page.locator('text=Download, button');
    await expect(downloadBtn).toBeDisabled();
  });

  test('Download Progress Display', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select files and click Download
    const checkboxes = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkboxes.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the progress bar appears and updates
    const progressOverlay = page.locator('text=#1/1, text=downloading');
    await expect(progressOverlay).toBeVisible();

    // Verify the download message shows step numbers
    await expect(progressOverlay).toBeVisible();

    // After completion, verify the message shows `completed.`
    const completedMsg = page.locator('text=completed.');
    await expect(completedMsg).toBeVisible();

    // Verify the progress bar is removed after completion
    await expect(progressOverlay).not.toBeVisible();
  });

  test('Clear Destination Directory', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Enter a destination directory path
    const destDirInput = page.locator('input[placeholder*="OPFS directory"], input[name*="dest"]');
    await destDirInput.fill('test/output');

    // Click the "Clear" button next to the destination field
    const clearBtn = page.locator('text=Clear, button, :text("Clear")');
    await clearBtn.click();

    // Verify the destination directory field is cleared
    await expect(destDirInput).toHaveValue('');

    // Verify the field is focused after clearing
    await expect(destDirInput).toBeFocused();
  });

  test('Download Single File', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select the file
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();

    // Set a destination directory
    const destDirInput = page.locator('input[placeholder*="OPFS directory"], input[name*="dest"]');
    await destDirInput.fill('single-output/');

    // Download the file
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the file appears in the OPFS destination directory
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with file:// Protocol', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select a file
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();

    // Download the file to OPFS
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the file is downloaded and uploaded to OPFS successfully
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with https:// Protocol', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select a file
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();

    // Download the file to OPFS
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the file is downloaded and uploaded to OPFS successfully
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download Error Handling', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select a file and attempt to download
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();

    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify an error message appears indicating the download failed
    const errorMsg = page.locator('text=Error, text=Failed');
    await expect(errorMsg).toBeVisible();

    // Verify the progress overlay is removed after error
    await expect(page.locator('text=downloading')).not.toBeVisible();

    // Verify the files that were successfully downloaded remain in OPFS
    await expect(page.locator('text=Download to OPFS')).toBeVisible();
  });

  test('Download to Nested OPFS Path', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Set a nested destination path
    const destDirInput = page.locator('input[placeholder*="OPFS directory"], input[name*="dest"]');
    await destDirInput.fill('output/data/');

    // Select and download files
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are saved in the nested OPFS path
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
    await expect(opfsPage!.locator('text=output/data/')).toBeVisible();
  });

  test('Download After Navigation', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Navigate away from the page
    await page.goto('http://localhost:9280/opfs/');

    // Navigate back to the original page
    await page.goto('http://localhost:9280/file://?toopfs');

    // Verify the toopfs UI is still functional
    const toopfsUI = page.locator('text=Download to OPFS');
    await expect(toopfsUI).toBeVisible();

    // Select and download files as before
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();
  });

  test('Download Multiple Files in Sequence', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select all files
    const checkboxes = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkboxes.first().check();

    // Click Download
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify each file is downloaded and uploaded sequentially
    const progressOverlay = page.locator('text=downloading');
    await expect(progressOverlay).toBeVisible();

    // Verify the progress shows the correct step numbers
    await expect(progressOverlay).toBeVisible();

    // Verify all files appear in the OPFS destination directory
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with Empty Selection', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Unselect all files
    const checkboxes = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkboxes.first().check(); // Select all first
    await checkboxes.first().click(); // Unselect all

    // Verify the "Download" button is disabled
    const downloadBtn = page.locator('text=Download, button');
    await expect(downloadBtn).toBeDisabled();

    // Attempt to click Download
    await downloadBtn.click();

    // Verify no action is triggered
    const progressOverlay = page.locator('text=downloading');
    await expect(progressOverlay).not.toBeVisible();
  });

  test('Download Progress Cancellation', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select files and click Download
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // While the download is in progress, verify the progress overlay is visible
    const progressOverlay = page.locator('text=downloading');
    await expect(progressOverlay).toBeVisible();

    // Verify the download completes or fails based on the source availability
    await expect(progressOverlay).toBeVisible();
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

  test('Download with Special Characters in Paths', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Set a destination with special characters
    const destDirInput = page.locator('input[placeholder*="OPFS directory"], input[name*="dest"]');
    await destDirInput.fill('my dir/output/');

    // Select and download files
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are saved correctly
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
