// spec: specs/opfs-ui-test-plan.md
// TC-04: ディレクトリへのナビゲーション
// TC-05: ブラウザの「戻る」ボタンによるナビゲーション

import { test, expect } from '@playwright/test';
import { gotoOPFS, createOPFSDir } from './helpers';

test.describe('TC-04: ディレクトリへのナビゲーション', () => {
  test('サブディレクトリに移動してパンくずで戻れる', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSDir(page, 'testdir');
    await page.click('#command-reload');

    // testdir/ リンクをクリックして移動
    await page.locator('.grid-row').filter({ hasText: 'testdir/' }).getByRole('link').click();

    // タイトルとパンくずが更新される
    await expect(page).toHaveTitle('OPFS: /testdir/');
    await expect(page.locator('#header')).toContainText('(Root)');
    await expect(page.locator('#header')).toContainText('testdir');

    // ルートの (Root) はリンクとして表示される
    await expect(page.locator('#header').getByRole('link', { name: '(Root)' })).toBeVisible();

    // (Root) リンクをクリックしてルートに戻る
    await page.locator('#header').getByRole('link', { name: '(Root)' }).click();

    await expect(page).toHaveTitle('OPFS: /');
    await expect(page.locator('#header')).toContainText('(Root)');
    // ルートでは (Root) はリンクではなく span になる
    await expect(page.locator('#header').getByRole('link')).toHaveCount(0);
  });
});

test.describe('TC-05: ブラウザの「戻る」ボタンによるナビゲーション', () => {
  test('ブラウザヒストリーと連動してディレクトリが切り替わる', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSDir(page, 'testdir');
    await page.click('#command-reload');

    // testdir に移動
    await page.locator('.grid-row').filter({ hasText: 'testdir/' }).getByRole('link').click();
    await expect(page).toHaveTitle('OPFS: /testdir/');

    // testdir 内に childdir を作成して移動
    await page.fill('#mkdir-name', 'childdir');
    await page.click('#mkdir-mkdir');
    await page.locator('.grid-row').filter({ hasText: 'childdir/' }).getByRole('link').click();
    await expect(page).toHaveTitle('OPFS: /testdir/childdir/');

    // ブラウザ「戻る」→ testdir に戻る
    await page.goBack();
    await expect(page).toHaveTitle('OPFS: /testdir/');

    // ブラウザ「戻る」→ ルートに戻る
    await page.goBack();
    await expect(page).toHaveTitle('OPFS: /');
  });
});
