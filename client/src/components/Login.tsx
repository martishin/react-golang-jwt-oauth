import React, { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Card from "./Card.tsx";

const Login: React.FC = () => {
  const navigate = useNavigate();

  useEffect(() => {
    /* global google */
    google.accounts.id.initialize({
      client_id: "YOUR_GOOGLE_APP_CLIENT_ID",
      callback: handleGoogleLogin,
    });

    google.accounts.id.renderButton(document.getElementById("google-signin-button")!, {
      theme: "outline",
      size: "large",
    });
  }, []);

  const handleGoogleLogin = async (response: any) => {
    try {
      const res = await fetch("http://localhost:8000/api/auth/google", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id_token: response.credential }),
        credentials: "include", // Ensure cookies are included
      });

      if (!res.ok) {
        throw new Error("Google login failed");
      }

      navigate("/secure");
    } catch (err: any) {
      console.error("Error!", err);
      alert("Error: " + err.message);
    }
  };

  return (
    <div className="flex items-center justify-center h-screen bg-gradient-to-br from-gray-100 to-gray-300">
      <Card>
        <h1 className="text-3xl font-bold text-gray-800 text-center mb-6">Welcome Back</h1>
        <p className="text-gray-600 text-center mb-6">
          Log in to continue using our service. It's quick and easy!
        </p>
        <div className="flex justify-center mb-4">
          <div id="google-signin-button"></div>
        </div>
        <div className="text-center mt-auto">
          <p className="text-sm text-gray-500">
            By logging in, you agree to our{" "}
            <a href="#" className="text-blue-500 hover:underline">
              Terms of Service
            </a>{" "}
            and{" "}
            <a href="#" className="text-blue-500 hover:underline">
              Privacy Policy
            </a>
            .
          </p>
        </div>
      </Card>
    </div>
  );
};

export default Login;
