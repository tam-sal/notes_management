import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axiosInstance from "../../config/axios";
import showToast from './../../context/toast-utils';
import { toast } from 'react-hot-toast';
import Loader from './../Loader/Loader';

const Note = ({ noteData, onArchive, onDelete }) => {
  const { title, content, is_archived: isArchived, categories, id } = noteData;
  const [loading, setLoading] = useState(false);
  const [localCategories, setLocalCategories] = useState(categories);
  const [isAddingCategory, setIsAddingCategory] = useState(false);
  const [categoryInput, setCategoryInput] = useState("");

  const navigate = useNavigate();

  const toggleArchiveStatus = async () => {
    toast.dismiss();
    setLoading(true);
    try {
      const { data } = await axiosInstance.put(`/notes/${id}/archive-toggle`);
      const updatedNote = data.data;
      showToast(true, `🟢 Successfully archived note with id: ${updatedNote?.id}`);
      onArchive(updatedNote.id, updatedNote.is_archived);
    } catch (err) {
      const error = `🔴 ${JSON.stringify(err.response?.data?.error || {})}`;
      showToast(false, error);
    } finally {
      setLoading(false);
    }
  };

  const deleteNote = async () => {
    toast.dismiss();
    setLoading(true);
    try {
      const { data } = await axiosInstance.delete(`/notes/${id}`);
      showToast(true, `🟢 Successfully deleted note with id: ${id}`);
      onDelete(id);
    } catch (err) {
      const error = `🔴 ${JSON.stringify(err.response?.data?.error || {})}`;
      showToast(false, error);
    } finally {
      setLoading(false);
    }
  };

  const addCategory = async () => {
    if (!categoryInput.trim()) {
      showToast(false, "🔴 Category name cannot be empty.");
      return;
    }
    setLoading(true);
    try {
      const { data } = await axiosInstance.post(`/notes/${id}/categories/${categoryInput.trim()}`);
      const updatedNote = data.data;
      showToast(true, `🟢 Successfully added category "${categoryInput.trim()}" to note.`);
      setLocalCategories(updatedNote.categories);
      setIsAddingCategory(false);
      setCategoryInput("");
    } catch (err) {
      const error = `🔴 ${JSON.stringify(err.response?.data?.error || {})}`;
      showToast(false, error);
    } finally {
      setLoading(false);
    }
  };

  const removeCategory = async (categoryName) => {
    setLoading(true);
    try {
      const { data } = await axiosInstance.delete(`/notes/${id}/categories/${categoryName}`);
      const updatedNote = data.data;
      showToast(true, `🟢 Successfully removed category "${categoryName}" from note.`);
      setLocalCategories(updatedNote.categories);
    } catch (err) {
      const error = `🔴 ${JSON.stringify(err.response?.data?.error || {})}`;
      showToast(false, error);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateClick = () => {
    navigate(`/notes/update/${id}`);
  };

  return (
    <>
      {loading && <Loader />}
      <div className="card bg-base-100 mb-3 pb-5 w-96 shadow-xl border-2 dark:border-slate-300 max-w-[85vw] mx-6">
        {/* Card Body */}
        <div className="card-body min-h-[60vh]">
          {/* Title */}
          <h2 className="card-title text-center">{title}</h2>
          {/* Content */}
          <p className='min-h-[40%] text-left'>{content}</p>
          {/* Categories */}
          <div className="card-actions justify-start flex-wrap min-h-[20%]">
            {localCategories?.map((category) => (
              <div key={category.id} className="badge badge-outline mr-2 mb-2">
                {category.name}
                <button
                  disabled={loading}
                  onClick={() => removeCategory(category.name)}
                  className="ml-2 text-red-500 hover:text-red-700"
                  aria-label={`Remove category ${category.name}`}
                >
                  ✕
                </button>
              </div>
            ))}
          </div>
          {!isAddingCategory ? (
            <button
              disabled={loading}
              onClick={() => setIsAddingCategory(true)}
              className="btn btn-sm btn-outline mt-4"
              aria-label="Add Category"
            >
              Add Category
            </button>
          ) : (
            <div className="flex items-center gap-2 mt-4">
              <input
                type="text"
                placeholder="Enter category name"
                className="input input-bordered w-full max-w-xs"
                value={categoryInput}
                onChange={(e) => setCategoryInput(e.target.value)}
                autoComplete="off"
              />
              <button
                disabled={loading}
                onClick={addCategory}
                className="btn btn-sm btn-primary"
                aria-label="Save Category"
              >
                Add
              </button>
            </div>
          )}
          {/* Actions */}
          <div className="card-actions justify-between mt-3 mb-1">
            {/* Update Note Button (Bottom Left) */}
            <button
              disabled={loading}
              onClick={handleUpdateClick}
              className="btn btn-sm btn-info"
              aria-label="Update Note"
            >
              Update
            </button>

            {/* Existing Actions (Bottom Right) */}
            <div className="flex gap-2">
              <button
                disabled={loading}
                onClick={toggleArchiveStatus}
                className={`btn btn-sm ${isArchived ? 'btn-success' : 'btn-warning'}`}
                aria-label={isArchived ? 'Unarchive Note' : 'Archive Note'}
              >
                {isArchived ? 'Unarchive' : 'Archive'}
              </button>
              <button
                disabled={loading}
                onClick={deleteNote}
                className="btn btn-sm btn-error"
                aria-label="Delete Note"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default Note;