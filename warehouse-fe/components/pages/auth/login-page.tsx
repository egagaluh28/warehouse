"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import {
  Input,
  Button,
  Card,
  CardHeader,
  CardBody,
  CardFooter,
  Image,
} from "@nextui-org/react";
import { EyeIcon, EyeSlashIcon } from "@heroicons/react/24/outline";
import { useAuth } from "@/components/providers/auth-provider";
import Link from "next/link";

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();
  const [formData, setFormData] = useState({
    username: "",
    password: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [isVisible, setIsVisible] = useState(false);

  const toggleVisibility = () => setIsVisible(!isVisible);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await login(formData);
      // Logic redirection moved to auth-provider or handled here if needed,
      // but auth-provider already handles it.
    } catch (err: any) {
      setError(
        err.response?.data?.message ||
          "Login failed. Please check your credentials.",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-white px-4 py-12 sm:px-6 lg:px-8">
      <div className="w-full max-w-md space-y-8">
        <Card className="p-4 shadow-lg">
          <CardHeader className="flex flex-col gap-1 pb-0 pt-4 text-center">
            <h2 className="text-2xl font-bold text-gray-900">
              ROXY Warehouse System
            </h2>
            <p className="text-sm text-gray-500">Sign in to your account</p>
          </CardHeader>
          <CardBody>
            <form onSubmit={handleSubmit} className="flex flex-col gap-4">
              <Input
                label="Username"
                placeholder="Enter your username"
                variant="bordered"
                value={formData.username}
                onValueChange={(val) =>
                  setFormData({ ...formData, username: val })
                }
                isRequired
              />
              <Input
                label="Password"
                placeholder="Enter your password"
                variant="bordered"
                value={formData.password}
                onValueChange={(val) =>
                  setFormData({ ...formData, password: val })
                }
                endContent={
                  <button
                    className="focus:outline-none"
                    type="button"
                    onClick={toggleVisibility}>
                    {isVisible ? (
                      <EyeSlashIcon className="h-5 w-5 text-gray-400" />
                    ) : (
                      <EyeIcon className="h-5 w-5 text-gray-400" />
                    )}
                  </button>
                }
                type={isVisible ? "text" : "password"}
                isRequired
              />
              {error && (
                <p className="text-center text-sm text-danger">{error}</p>
              )}
              <Button
                color="primary"
                type="submit"
                isLoading={loading}
                className="w-full font-semibold text-white">
                Sign in
              </Button>
            </form>
          </CardBody>
          <CardFooter className="justify-center">
            <p className="text-sm text-gray-500">
              Don't have an account? Contact Admin.
            </p>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
