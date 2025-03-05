import { useState, useEffect } from 'react';
import axiosInstance from "../../config/axios";
import showToast from "../../context/toast-utils";
import { toast } from 'react-hot-toast';
import Loader from '../Loader/Loader';
import { useParams, useNavigate } from 'react-router-dom';

function UpdateNote() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [formData, setFormData] = useState({
    title: "",
    content: "",
    categories: "",
  });
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [loading, setLoading] = useState(false);


  useEffect(() => {
    const fetchNote = async () => {
      setLoading(true);
      try {
        const response = await axiosInstance.get(`/notes/${id}`);
        const note = response.data.data;
        const categoriesString = note.categories.map(cat => cat.name).join(", ");

        setFormData({
          title: note.title,
          content: note.content,
          categories: categoriesString,
        });
      } catch (err) {
        const error = `🔴 ${Object.values(err.response?.data?.error || {}).join(' ')}` || "🔴 Failed to fetch note details.";
        setError(error);
        setTimeout(() => showToast(false, error), 200);
      } finally {
        setLoading(false);
      }
    };
    fetchNote();
  }, [id]);

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

    const categoriesArray = formData.categories
      .split(",")
      .map(cat => ({ name: cat.trim() }));

    try {
      const response = await axiosInstance.put(`/notes/${id}`, {
        title: formData.title,
        content: formData.content,
        categories: categoriesArray,
      });

      const successMessage = `🟢 Note updated successfully with ID: ${id}`;
      setSuccess(successMessage);
      showToast(true, successMessage);
      setTimeout(() => navigate("/notes"), 100);
    } catch (err) {
      const errorMessage =
        `🔴 ${Object.values(err.response?.data?.error || {}).join(' ')}` ||
        "🔴 Something went wrong. Try again.";
      setError(errorMessage);
      showToast(false, errorMessage)
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      {loading && <Loader />}
      <div className="min-w-[30vw] mx-auto p-4 border rounded-lg shadow bg-base-100 text-base-content max-w-[90vw]">
        <h2 className="text-xl font-bold mb-4 text-center">Update Note</h2>
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
            Update Note
          </button>
        </form>
      </div>
    </>
  );
}

export default UpdateNote;