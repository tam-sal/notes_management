import { React, useState } from 'react'
import axiosInstance from "../../config/axios"
import showToast from "../../context/toast-utils";
import { toast } from 'react-hot-toast';
import Loader from '../Loader/Loader';
import { useNavigate } from 'react-router-dom';

function CreateNote() {

  const navigate = useNavigate()
  const initialForm = {
    title: "",
    content: "",
    categories: "",
  };
  let [formData, setFormData] = useState({
    title: "",
    content: "",
    categories: "",
  });

  let [error, setError] = useState("");
  let [success, setSuccess] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleSubmit = async (e) => {
    toast.dismiss();
    setLoading(true);
    e.preventDefault();
    setError("");
    setSuccess("");

    const categoriesArray = formData.categories.split(",").map((cat) => cat.trim());

    try {
      const response = await axiosInstance.post("/notes", {
        title: formData.title,
        content: formData.content,
        categories: categoriesArray,
      });
      success = `🟢 Note created with ID: ${response.data.data.id}`
      setSuccess(success);
      setFormData(initialForm)
      setTimeout(() => navigate("/notes"), 100);
      setTimeout(() => showToast(true, success), 200);
    } catch (err) {
      error = `🔴 ${Object.values(err.response?.data?.error).join(' ')}` || "🔴 Something went wrong. Try again."
      setError(error);
      setTimeout(() => showToast(false, error), 200);
    } finally {
      toast.dismiss();
      setLoading(false);
    }
  };

  return (
    <>
      {loading && <Loader />}
      <div className="min-w-[30vw] mx-auto p-4 border rounded-lg shadow bg-base-100 text-base-content max-w-[90vw]">
        <h2 className="text-xl font-bold mb-4 text-center">Create Note</h2>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium mb-1" htmlFor="title">
              Title
            </label>
            <input
              id="title"
              name="title"
              type="text"
              value={formData.title}
              onChange={handleChange}
              className="input input-bordered w-full bg-base-200 text-base-content placeholder-base-content/60 focus:outline-none focus:border-primary"
              required
              autoComplete="off"
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium mb-1" htmlFor="content">
              Content
            </label>
            <textarea
              id="content"
              name="content"
              rows="4"
              value={formData.content}
              onChange={handleChange}
              className="textarea textarea-bordered w-full bg-base-200 text-base-content placeholder-base-content/60 focus:outline-none focus:border-primary"
              required
            ></textarea>
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium mb-1" htmlFor="categories">
              Categories (comma-separated)
            </label>
            <input
              id="categories"
              name="categories"
              type="text"
              value={formData.categories}
              onChange={handleChange}
              className="input input-bordered w-full bg-base-200 text-base-content placeholder-base-content/60 focus:outline-none focus:border-primary"
              placeholder="e.g., science fiction, favorites"
              autoComplete="off"
            />
          </div>
          <button
            disabled={loading}
            type="submit"
            className="btn btn-primary w-full"
          >
            Create Note
          </button>
        </form>
      </div>
    </>
  );
}

export default CreateNote