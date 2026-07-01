import { test, expect } from '@playwright/test';

test.describe('01-navigation', () => {
  test('OPFSページが正常に読み込まれる', async ({ page }) => {
    // 1. http://127.0.0.1:9280/opfs/ にアクセス
    await page.goto('http://127.0.0.1:9280/opfs/');

    // expect: タイトルが"OPFS: /"である
    await expect(page).toHaveTitle('OPFS: /');

    // expect: h1に"OPFS: Origin Private File System"が表示される
    await expect(page.getByRole('heading', { level: 1 })).toBeVisible();
    await expect(page.getByRole('heading', { level: 1 })).toHaveText(/OPFS: Origin Private File System/);

    // expect: パンくずリストに"(Root)"が表示される
    await expect(page.locator('#header span')).toBeVisible();
  });

  test('ページタイトルを確認', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // expect: ページタイトルが"OPFS: /"である
    await expect(page).toHaveTitle('OPFS: /');
  });

  test('h1要素の内容を確認', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // expect: h1要素に"OPFS: Origin Private File System"が含まれる
    const h1 = page.getByRole('heading', { level: 1 });
    await expect(h1).toBeVisible();
    const text = await h1.textContent();
    expect(text).toContain('OPFS: Origin Private File System');
  });

  test('パンくずリストの内容を確認', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // expect: パンくずリストに"(Root)"が表示される
    await expect(page.locator('#header span')).toBeVisible();
  });

  test('サブディレクトリのナビゲーション（クリック）', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // 1. ディレクトリ作成: #mkdir-nameに"nav-test"を入力し、#mkdir-mkdirをクリック
    await page.locator('#mkdir-name').fill('nav-test');
    await page.locator('#mkdir-mkdir').click();

    // expect: サブディレクトリが作成されている
    await expect(page.locator('a:has-text("nav-test/")')).toBeVisible();

    // 2. "nav-test/"のリンクをクリック
    await page.locator('a:has-text("nav-test/")').click();

    // expect: ページタイトルが"OPFS: /nav-test/"である
    await expect(page).toHaveTitle('OPFS: /nav-test/');

    // 3. location.hashが"#/nav-test/"であることを確認
    const hash = await page.evaluate(() => window.location.hash);
    expect(hash).toBe('#/nav-test/');

    // 4. ページタイトルを確認
    await expect(page).toHaveTitle('OPFS: /nav-test/');

    // 5. パンくずリストを確認
    await expect(page.locator('#header span')).toBeVisible();
  });

  test('パンくずリストでの親ディレクトリ移動', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create nav-test directory first
    await page.locator('#mkdir-name').fill('nav-test');
    await page.locator('#mkdir-mkdir').click();
    await page.locator('a:has-text("nav-test/")').click();
    await expect(page).toHaveTitle('OPFS: /nav-test/');

    // 1. サブディレクトリ内にいる状態で、パンくずの"(Root)"をクリック
    await page.locator('#header span').click();

    // expect: パンくずの"(Root)"リンクをクリックするとルートに戻る

    // 2. ページタイトルが"OPFS: /"であることを確認
    await expect(page).toHaveTitle('OPFS: /');

    // 3. location.hashを確認
    const hash = await page.evaluate(() => window.location.hash);
    expect(hash).toBe('#');
  });

  test('操作後のリレンダリング', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // 1. ディレクトリを作成し、一覧に追加されることを確認
    await page.locator('#mkdir-name').fill('render-test-dir');
    await page.locator('#mkdir-mkdir').click();
    await expect(page.locator('a:has-text("render-test-dir/")')).toBeVisible();

    // 2. エディタでファイルを作成し、一覧に追加されることを確認
    await page.locator('#editor-name').fill('render-test.txt');
    await page.locator('#editor-edit').fill('Hello, OPFS!');
    await page.locator('#editor-save').click();

    // expect: ファイル作成後、一覧に新ファイルが表示される
    await expect(page.locator('a:has-text("render-test.txt")')).toBeVisible();

    // 3. ファイルを削除し、一覧から消えることを確認
    await page.locator('input[name="render-test.txt"]').check();
    await page.locator('#command-delete').click();

    // Handle confirm dialog - accept
    page.once('dialog', dialog => {
      expect(dialog.message()).toContain('render-test.txt');
      dialog.accept();
    });

    await expect(page.locator('a:has-text("render-test.txt")')).not.toBeVisible();
  });
});
