import { create } from "zustand";

/**
 * 認証フロー中の一時状態。
 * register / login で OTP を送った宛先メールを保持し、OTP 検証画面で参照する。
 * 機密情報(トークン等)は保持しない（保管方式は nextjs-best-practices.md の論点が未確定）。
 */
type AuthFlowState = {
	pendingEmail: string | null;
	setPendingEmail: (email: string | null) => void;
};

export const useAuthFlowStore = create<AuthFlowState>()((set) => ({
	pendingEmail: null,
	setPendingEmail: (pendingEmail) => set({ pendingEmail }),
}));
