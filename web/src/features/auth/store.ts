import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

/**
 * 認証フロー中の一時状態。
 * register / login で OTP を送った宛先メールを保持し、OTP 検証画面で参照する。
 * 機密情報(トークン等)は保持しない(トークンは httpOnly Cookie 側)。
 *
 * `sessionStorage` に永続化する理由:
 * - OTP 画面はメールアプリへ切り替える導線が多く、リロード/直アクセスで
 *   メモリ上の状態が消えるとフロー続行不能になるため。
 * - `sessionStorage` はタブを閉じれば破棄され、URL にもメールを残さない
 *   (一時状態としての性質に合致)。
 *
 * `hasHydrated` は `sessionStorage` からの復元完了フラグ。復元前は `pendingEmail`
 * が常に `null` のため、これを待たずにリダイレクト判定するとフラッシュが起きる。
 */
type AuthFlowState = {
	pendingEmail: string | null;
	setPendingEmail: (email: string | null) => void;
	hasHydrated: boolean;
	setHasHydrated: (value: boolean) => void;
};

export const useAuthFlowStore = create<AuthFlowState>()(
	persist(
		(set) => ({
			pendingEmail: null,
			setPendingEmail: (pendingEmail) => set({ pendingEmail }),
			hasHydrated: false,
			setHasHydrated: (hasHydrated) => set({ hasHydrated }),
		}),
		{
			name: "spark-auth-flow",
			// SSR では sessionStorage 未定義 → createJSONStorage が握りつぶし no-op 化する。
			storage: createJSONStorage(() => sessionStorage),
			// 永続化対象は宛先メールのみ(復元完了フラグは毎回 false から始める)。
			partialize: (state) => ({ pendingEmail: state.pendingEmail }),
			onRehydrateStorage: () => (state) => {
				state?.setHasHydrated(true);
			},
		},
	),
);
