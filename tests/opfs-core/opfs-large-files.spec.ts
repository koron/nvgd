import { test, expect } from '@playwright/test';

test.describe('Edge Cases and Error Handling', () => {
  test('OPFS Operations with Very Large Files', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a file with content larger than 1MB in OPFS
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    const editorTextarea = page.locator('textarea');
    const largeContent = 'A'.repeat(1024 * 1024); // 1MB
    await editorName.fill('large-file.txt');
    await editorTextarea.fill(largeContent);
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "large-file.txt"')).toBeVisible();

    // Step 2: Verify the file appears in the listing with correct size
    const listing = page.locator('.listing, .directory-list, table');
    await expect(listing).toContainText('large-file.txt');

    // Step 3: Verify the file can be downloaded via "Save as" action
    const saveAsBtn = page.locator('text=Save as, [class*="save"], button[title*="Save"]');
    await expect(saveAsBtn).toBeVisible();

    // Step 4: Verify the file can be uploaded via the file picker
    const uploadInput = page.locator('input[type=file]');
    await uploadInput.setInputFiles('tests/testdata/sample.txt');
    await page.locator('text=Upload, button').click();
    await expect(page.locator('text=Uploaded "sample.txt"')).toBeVisible();
  });

  test('OPFS Operations with Unicode Filenames', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a file with a Unicode name (e.g., `日本語.txt`)
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('日本語.txt');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('Unicode filename content');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=日本語.txt')).toBeVisible();

    // Step 2: Verify the file appears correctly in the listing
    await expect(page.locator('text=日本語.txt')).toBeVisible();

    // Step 3: Upload a file with a Unicode name
    const uploadInput = page.locator('input[type=file]');
    await uploadInput.setInputFiles('tests/testdata/sample.txt');
    const nameField = page.locator('input[placeholder*="File name"]');
    await nameField.fill('файл.txt');
    await page.locator('text=Upload, button').click();
    await expect(page.locator('text=файл.txt')).toBeVisible();

    // Step 4: Verify the file appears correctly in the listing
    await expect(page.locator('text=файл.txt')).toBeVisible();
  });

  test('OPFS Operations with Very Long Filenames', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a file with a very long name (e.g., 255 characters)
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    const veryLongName = 'a'.repeat(255) + '.txt';
    await editorName.fill(veryLongName);
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('Long filename content');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded')).toBeVisible();

    // Step 2: Verify the file appears correctly in the listing
    const listing = page.locator('.listing, .directory-list, table');
    await expect(listing).toContainText(veryLongName);

    // Step 3: Verify the file can be downloaded via "Save as" action
    const saveAsBtn = page.locator('text=Save as, [class*="save"], button[title*="Save"]');
    await expect(saveAsBtn).toBeVisible();
  });

  test('OPFS Operations with Special Characters in Paths', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a directory with a name containing special characters (e.g., `my dir/`)
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('my dir');
    await page.locator('text=Create new directory, button').click();
    await expect(page.locator('text=my dir/')).toBeVisible();

    // Step 2: Create a file inside the directory
    await page.locator('text=my dir/').click();
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('special.txt');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('Special characters in path');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "special.txt"')).toBeVisible();

    // Step 3: Verify the file appears correctly in the listing
    await expect(page.locator('text=special.txt')).toBeVisible();

    // Step 4: Verify the breadcrumb correctly displays the directory with special characters
    await expect(page.locator('text=my dir')).toBeVisible();
  });

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

  test('OPFS Operations with Network Errors', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Navigate to an NVGD page with a broken source URL and `?toopfs` filter
    const urlInput = page.locator('input[placeholder*="URL"]');
    await urlInput.fill('http://localhost:9280/nonexistent-page');
    const nameInput = page.locator('input[placeholder*="Name"]');
    await nameInput.fill('broken.md');

    // Step 2: Select a file and attempt to download
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();

    // Step 3: Verify an error message appears indicating the download failed
    const errorMsg = page.locator('text=Error, text=Failed');
    await expect(errorMsg).toBeVisible();

    // Step 4: Verify the progress overlay is removed after error
    await expect(page.locator('text=downloading')).not.toBeVisible();
  });

  test('OPFS Operations with Timeout', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Navigate to an NVGD page with a slow source URL and `?toopfs` filter
    const urlInput = page.locator('input[placeholder*="URL"]');
    await urlInput.fill('http://localhost:9280/help/');
    const nameInput = page.locator('input[placeholder*="Name"]');
    await nameInput.fill('timeout-test.md');

    // Step 2: Select a file and attempt to download
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();

    // Step 3: Verify the download completes or times out
    const errorMsg = page.locator('text=Error, text=Failed');
    await expect(errorMsg).toBeVisible();

    // Step 4: Verify the progress overlay is removed after timeout
    await expect(page.locator('text=downloading')).not.toBeVisible();
  });

  test('OPFS Operations with Invalid Input', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Enter an invalid directory name in the "Directory name to create" input
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('');
    await page.locator('text=Create new directory, button').click();
    const dirError = page.locator('text=Need directory name');
    await expect(dirError).toBeVisible();

    // Step 2: Enter an invalid file name in the "File name to create or update" field
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('Invalid file name');
    await page.locator('text=Create or update a file, button').click();
    const fileError = page.locator('text=Need file name');
    await expect(fileError).toBeVisible();

    // Step 3: Enter an invalid URL in the "URL to be downloaded" field
    const urlInput = page.locator('input[placeholder*="URL"]');
    await urlInput.fill('not-a-valid-url');
    const nameInput = page.locator('input[placeholder*="Name"]');
    await nameInput.fill('invalid.md');
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    const urlError = page.locator('text=Error, text=Failed');
    await expect(urlError).toBeVisible();
  });
});
