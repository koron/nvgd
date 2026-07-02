import { test, expect } from '@playwright/test';

test.describe('DuckDB Integration', () => {
  test('Open Supported File Types', async ({ page, browser }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a CSV file in OPFS with tabular data
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('data.csv');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('name,age,city\nAlice,30,NYC\nBob,25,LA\nCharlie,35,Chicago');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "data.csv"')).toBeVisible();

    // Step 2: Check the checkbox next to the CSV file
    const csvCheckbox = page.locator('text=data.csv').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await csvCheckbox.check();

    // Step 3: Click the "DuckDB" button
    const duckdbBtn = page.locator('text=DuckDB, button');
    await expect(duckdbBtn).toBeEnabled();
    await duckdbBtn.click();

    // Step 4: Verify a DuckDB WASM shell opens in a new tab
    await page.waitForEvent('page');
    const duckdbPage = page.context().pages().find(p => p.url().includes('duckdb'));
    expect(duckdbPage).toBeTruthy();

    // Step 5: Verify the DuckDB shell shows a view named `opfs0`
    const viewOpfs0 = duckdbPage!.locator('text=opfs0, .view, [class*="view"]');
    await expect(viewOpfs0).toBeVisible();

    // Step 6: Verify the DuckDB shell shows the data from the CSV file
    const dataDisplay = duckdbPage!.locator('text=Alice, Bob, Charlie, .result, .output');
    await expect(dataDisplay).toBeVisible();

    // Step 7: Close the DuckDB tab
    await duckdbPage!.close();
  });

  test('Open Multiple Files', async ({ page, browser }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create two CSV files in OPFS
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    const editorTextarea = page.locator('textarea');

    await editorName.fill('data1.csv');
    await editorTextarea.fill('name,value\nA,10\nB,20');
    await page.locator('text=Create or update a file, button').click();

    await editorName.fill('data2.csv');
    await editorTextarea.fill('name,value\nC,30\nD,40');
    await page.locator('text=Create or update a file, button').click();

    // Step 2: Check both file checkboxes
    const checkbox1 = page.locator('text=data1.csv').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await checkbox1.check();
    const checkbox2 = page.locator('text=data2.csv').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await checkbox2.check();

    // Step 3: Click the "DuckDB" button
    const duckdbBtn = page.locator('text=DuckDB, button');
    await duckdbBtn.click();

    // Step 4: Verify the DuckDB shell shows views `opfs0` and `opfs1`
    await page.waitForEvent('page');
    const duckdbPage = page.context().pages().find(p => p.url().includes('duckdb'));
    expect(duckdbPage).toBeTruthy();

    const viewOpfs0 = duckdbPage!.locator('text=opfs0');
    await expect(viewOpfs0).toBeVisible();
    const viewOpfs1 = duckdbPage!.locator('text=opfs1');
    await expect(viewOpfs1).toBeVisible();

    // Step 5: Verify each view contains data from the corresponding file
    await expect(duckdbPage!.locator('text=10, 20')).toBeVisible();
    await expect(duckdbPage!.locator('text=30, 40')).toBeVisible();

    // Close DuckDB tab
    await duckdbPage!.close();
  });

  test('Open Directory Recursively', async ({ page, browser }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a directory with multiple CSV files inside
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('csvdir');
    await page.locator('text=Create new directory, button').click();

    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    const editorTextarea = page.locator('textarea');

    await editorName.fill('file1.csv');
    await editorTextarea.fill('col1,col2\n1,2\n3,4');
    await page.locator('text=Create or update a file, button').click();

    await editorName.fill('file2.csv');
    await editorTextarea.fill('col1,col2\n5,6\n7,8');
    await page.locator('text=Create or update a file, button').click();

    // Step 2: Check the checkbox next to the directory
    const dirCheckbox = page.locator('text=csvdir/').locator('input[type=checkbox], [class*="checkbox"], [class*="check"]');
    await dirCheckbox.check();

    // Step 3: Click the "DuckDB" button
    const duckdbBtn = page.locator('text=DuckDB, button');
    await duckdbBtn.click();

    // Step 4: Verify the DuckDB shell creates views for all files in the directory
    await page.waitForEvent('page');
    const duckdbPage = page.context().pages().find(p => p.url().includes('duckdb'));
    expect(duckdbPage).toBeTruthy();

    const viewOpfs0 = duckdbPage!.locator('text=opfs0');
    await expect(viewOpfs0).toBeVisible();
    const viewOpfs1 = duckdbPage!.locator('text=opfs1');
    await expect(viewOpfs1).toBeVisible();

    // Step 5: Verify all views are queryable
    await expect(duckdbPage!.locator('text=1, 2')).toBeVisible();
    await expect(duckdbPage!.locator('text=5, 6')).toBeVisible();

    // Close DuckDB tab
    await duckdbPage!.close();
  });

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
