import { OtpForm } from "@/features/auth/components/otp-form";

/**
 * OTP 検証画面 (`/otp`)。Server Component。
 * 対話的なフォーム本体は末端の Client Component (`OtpForm`) に委譲する。
 */
export default function OtpPage() {
	return <OtpForm />;
}
