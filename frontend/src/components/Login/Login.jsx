import { Link, useNavigate } from "react-router-dom";
import { useState } from "react";
import axiosInstance from "../../config/axios";
import showToast from "../../context/toast-utils";
import { toast } from "react-hot-toast";
import { useAuth } from "../../context/AuthContext"; // Import useAuth
import Loader from "../Loader/Loader";

const Login = () => {
  const initialForm = { user_name: "", password: "" };
  const [form, setForm] = useState(initialForm);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { checkAuth } = useAuth();

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    toast.dismiss();
    setLoading(true);

    try {
      const { data } = await axiosInstance.post("/user/login", form);
      const successMsg = `${data.data?.user_name} successfully logged in`;
      await checkAuth().catch(error => {
        console.error("AUTH CHECK FAILED:", error);
        throw new Error("AUTH FAILED")
      });
      showToast(true, successMsg);
      navigate("/notes");
    } catch (err) {
      const errMsg =
        `🔴 ${Object.values(err.response?.data?.error).join(" ")}` ||
        "🔴 Something went wrong. Try again.";
      showToast(false, errMsg);
    } finally {
      setLoading(false);
      setForm(initialForm);
    }
  };

  return (
    <>
      {loading && <Loader />}
      <div className="login min-w-[20vw]">
        <form
          className="max-w-md mx-auto p-8 bg-base-100 rounded-box shadow-lg border border-primary/20 min-w-full"
          onSubmit={handleSubmit}
        >
          {/* Title */}
          <h2 className="text-2xl font-bold mb-4 text-base-content">Log In</h2>

          {/* Username Field */}
          <div className="mb-6">
            <label className="flex items-center gap-2 label">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 16 16"
                fill="currentColor"
                className="h-4 w-4 opacity-70 text-inherit"
              >
                <path d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM12.735 14c.618 0 1.093-.561.872-1.139a6.002 6.002 0 0 0-11.215 0c-.22.578.254 1.139.872 1.139h9.47Z" />
              </svg>
              <input
                type="text"
                id="user_name"
                name="user_name"
                className="input input-bordered input-primary grow"
                placeholder="Username"
                value={form.user_name}
                onChange={handleInputChange}
                required
                autoComplete="off"
              />
            </label>
          </div>

          {/* Password Field */}
          <div className="mb-6">
            <label className="flex items-center gap-2 label">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 16 16"
                fill="currentColor"
                className="h-4 w-4 opacity-70 text-inherit"
              >
                <path
                  fillRule="evenodd"
                  d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
                  clipRule="evenodd"
                />
              </svg>
              <input
                type="password"
                id="password"
                name="password"
                className="input input-bordered input-primary grow"
                placeholder="Password"
                value={form.password}
                onChange={handleInputChange}
                required
                autoComplete="off"
              />
            </label>
          </div>

          {/* Submit Section */}
          <div className="flex items-center justify-between mt-8">
            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading}
            >
              {loading ? <span className="loading loading-spinner"></span> : "Log in"}
            </button>
            <p className="text-sm link link-primary">
              <Link to="/register">Don't have an account?</Link>
            </p>
          </div>
        </form>
      </div>
    </>
  );
};

export default Login;