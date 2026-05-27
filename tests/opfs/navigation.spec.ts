// spec: specs/opfs-ui-test-plan.md
// TC-04: ディレクトリへのナビゲーション
// TC-05: ブラウザの「戻る」ボタンによるナビゲーション

import { test, expect } from '@playwright/test';
import { gotoOPFS, createOPFSDir } from './helpers';

const SKIP_WEBKIT = 'OPFS は WebKit の HTTP では利用不可（セキュアコンテキスト外）';

test.describe('TC-04: ディレクトリへのナビゲーション', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('サブディレクトリに移動してパンくずで戻れる', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSDir(page, 'testdir');
    await page.click('#command-reload');

    // testdir/ リンクをクリックして移動
    // Mithril は href なしの <a> をレンダリングするため getByRole('link') ではなく locator('a') を使う
    await page.locator('.grid-row').filter({ hasText: 'testdir/' }).locator('a').click();

    // タイトルとパンくずが更新される
    await expect(page).toHaveTitle('OPFS: /testdir/');
    await expect(page.locator('#header')).toContainText('(Root)');
    await expect(page.locator('#header')).toContainText('testdir');

    // ルートの (Root) はクリッカブルな <a> として表示される
    await expect(page.locator('#header').locator('a', { hasText: '(Root)' })).toBeVisible();

    // (Root) リンクをクリックしてルートに戻る
    await page.locator('#header').locator('a', { hasText: '(Root)' }).click();

    await expect(page).toHaveTitle('OPFS: /');
    await expect(page.locator('#header')).toContainText('(Root)');
    // ルートでは (Root) は span になりリンクが消える
    await expect(page.locator('#header').locator('a')).toHaveCount(0);
  });
});

test.describe('TC-05: ブラウザの「戻る」ボタンによるナビゲーション', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('ブラウザヒストリーと連動してディレクトリが切り替わる', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSDir(page, 'testdir');
    await page.click('#command-reload');

    // testdir に移動
    await page.locator('.grid-row').filter({ hasText: 'testdir/' }).locator('a').click();
    await expect(page).toHaveTitle('OPFS: /testdir/');

    // testdir 内に childdir を作成して移動
    await page.fill('#mkdir-name', 'childdir');
    await page.click('#mkdir-mkdir');
    await page.locator('.grid-row').filter({ hasText: 'childdir/' }).locator('a').click();
    await expect(page).toHaveTitle('OPFS: /testdir/childdir/');

    // ブラウザ「戻る」→ testdir に戻る
    await page.goBack();
    await expect(page).toHaveTitle('OPFS: /testdir/');

    // ブラウザ「戻る」→ ルートに戻る
    await page.goBack();
    await expect(page).toHaveTitle('OPFS: /');
  });
});
