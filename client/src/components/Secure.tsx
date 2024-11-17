import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Card from "./Card.tsx";

interface UserDetails {
  name: string;
  email: string;
  picture?: string;
}

const Secure: React.FC = () => {
  const navigate = useNavigate();
  const [userDetails, setUserDetails] = useState<UserDetails | null>(null);

  const fetchUserDetails = async () => {
    try {
      const res = await fetch("http://localhost:8000/api/user", {
        credentials: "include",
      });

      if (res.status === 401) {
        const refreshRes = await fetch("http://localhost:8000/api/auth/refresh", {
          method: "POST",
          credentials: "include",
        });

        if (!refreshRes.ok) {
          throw new Error("Failed to refresh token");
        }

        const retryRes = await fetch("http://localhost:8000/api/user", {
          credentials: "include",
        });

        if (!retryRes.ok) {
          throw new Error("Failed to fetch user");
        }

        const userData = await retryRes.json();
        setUserDetails(userData);
      } else if (res.ok) {
        const userData = await res.json();
        setUserDetails(userData);
      } else {
        throw new Error("Failed to fetch user");
      }
    } catch (err) {
      console.error("Error fetching user:", err);
      navigate("/");
    }
  };

  useEffect(() => {
    fetchUserDetails();

    const intervalId = setInterval(
      () => {
        fetch("http://localhost:8000/api/auth/refresh", {
          method: "POST",
          credentials: "include",
        }).catch((err) => console.error("Error refreshing token:", err));
      },
      10 * 60 * 1000,
    );

    return () => clearInterval(intervalId);
  }, [navigate]);

  const handleLogout = async () => {
    try {
      await fetch("http://localhost:8000/api/auth/logout", {
        method: "POST",
        credentials: "include",
      });
      setUserDetails(null);
      navigate("/");
    } catch (err) {
      console.error("Error during logout:", err);
      alert("Logout failed. Please try again.");
    }
  };

  return (
    <>
      {userDetails ? (
        <div className="flex items-center justify-center h-screen w-screen bg-gray-100 bg-gradient-to-br from-gray-100 to-gray-200 font-sans">
          <Card>
            {userDetails.picture && (
              <img
                src={userDetails.picture}
                alt={`${userDetails.name}'s profile`}
                className="w-30 h-30 rounded-full mb-5 border-4 border-gray-200 mx-auto"
              />
            )}
            <p>Welcome</p>
            <h1 className="text-gray-800 my-3 text-3xl font-semibold">{userDetails.name}</h1>
            <p className="text-gray-600 text-base my-2">{userDetails.email}</p>
            <div className="flex justify-center mt-auto">
              <button
                onClick={handleLogout}
                className="bg-red-600 text-white py-2 px-4 rounded hover:bg-red-500"
              >
                Logout
              </button>
            </div>
          </Card>
        </div>
      ) : (
        <div className="flex items-center justify-center h-screen w-screen">
          <h1>Loading...</h1>
        </div>
      )}
    </>
  );
};

export default Secure;
