import React, { useState, useEffect } from 'react';
import axiosInstance from '../../config/axios';
import Note from './Note';
import { toast } from 'react-hot-toast';
import Loader from '../Loader/Loader';

const Notes = () => {
  const [notes, setNotes] = useState([]);
  const [activeFilter, setActiveFilter] = useState('all');
  const [loading, setLoading] = useState(true);
  const [errorMessage, setErrorMessage] = useState("");
  const [selectedCategories, setSelectedCategories] = useState([]);
  const [categoryInput, setCategoryInput] = useState("");

  // Function to fetch notes based on the filter
  const fetchNotes = async () => {
    setLoading(true);
    setErrorMessage(null);
    try {
      let url = '/notes';

      // Determine if we need to use /filter
      const hasFilters = activeFilter !== 'all' || selectedCategories.length > 0;
      if (hasFilters) {
        url = '/notes/filter';
      }

      const queryParams = [];
      if (activeFilter === 'active') {
        queryParams.push('isArchived=false');
      } else if (activeFilter === 'archived') {
        queryParams.push('isArchived=true');
      }

      if (selectedCategories.length > 0) {
        selectedCategories.forEach(cat => {
          queryParams.push(`categories=${encodeURIComponent(cat)}`);
        });
      }

      // Append query parameters to the URL
      if (queryParams.length > 0) {
        url += `?${queryParams.join('&')}`;
      }

      const response = await axiosInstance.get(url);

      if (!response.data.data || response.data.data.length === 0) {
        setNotes([]);
      } else {
        setNotes(response.data.data);
      }
    } catch (err) {
      setErrorMessage(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNotes();
    return () => {
      setNotes([]);
    };
  }, [activeFilter, selectedCategories]);

  const handleFilterClick = (filter) => {
    setActiveFilter(filter);
  };

  const addCategory = () => {
    const trimmedCategory = categoryInput.trim();
    if (trimmedCategory && !selectedCategories.includes(trimmedCategory)) {
      setSelectedCategories([...selectedCategories, trimmedCategory]);
      setCategoryInput("");
    }
  };

  const removeCategory = (categoryToRemove) => {
    setSelectedCategories(selectedCategories.filter(cat => cat !== categoryToRemove));
  };

  return (
    <div className="p-4">
      {/* Buttons for filtering */}
      <div className="flex gap-2 mb-4 flex-wrap">
        <button
          disabled={loading}
          className={`btn ${activeFilter === 'all' ? 'btn-primary' : 'btn-outline'}`}
          onClick={() => handleFilterClick('all')}
        >
          All Notes
        </button>
        <button
          disabled={loading}
          className={`btn ${activeFilter === 'active' ? 'btn-primary' : 'btn-outline'}`}
          onClick={() => handleFilterClick('active')}
        >
          Active Notes
        </button>
        <button
          disabled={loading}
          className={`btn ${activeFilter === 'archived' ? 'btn-primary' : 'btn-outline'}`}
          onClick={() => handleFilterClick('archived')}
        >
          Archived Notes
        </button>

        {/* Category Filter Input */}
        <div className="flex items-center gap-2">
          <input
            type="text"
            placeholder="Add category filter"
            className="input input-bordered w-full max-w-xs"
            value={categoryInput}
            onChange={(e) => setCategoryInput(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && addCategory()}
            autoComplete="off"
          />
          <button
            disabled={loading}
            onClick={addCategory}
            className="btn btn-sm btn-primary"
            aria-label="Add Category Filter"
          >
            Add
          </button>
        </div>

        {/* Selected Categories */}
        <div className="flex gap-2 flex-wrap mt-2">
          {selectedCategories.map((category) => (
            <div key={category} className="badge badge-outline">
              {category}
              <button
                disabled={loading}
                onClick={() => removeCategory(category)}
                className="ml-2 text-red-500 hover:text-red-700"
                aria-label={`Remove category filter ${category}`}
              >
                ✕
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Conditional rendering based on loading, error, and notes */}
      {loading && <Loader />}
      {!loading && notes.length === 0 && (
        <h2 className="text-center text-lg font-bold">No Notes Were Found</h2>
      )}
      {!loading && !errorMessage && notes.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {notes.map((note) => (
            <Note
              key={note.id}
              noteData={note}
              onArchive={(id, isArchived) => {
                // Update the note in the list after archiving
                setNotes((prevNotes) => {
                  const updatedNotes = prevNotes.map((n) =>
                    n.id === id ? { ...n, is_archived: isArchived } : n
                  );

                  // Filter the notes based on the active filter
                  if (activeFilter === 'active' && isArchived) {
                    // Remove archived notes from the "Active Notes" view
                    return updatedNotes.filter((n) => !n.is_archived);
                  } else if (activeFilter === 'archived' && !isArchived) {
                    // Remove unarchived notes from the "Archived Notes" view
                    return updatedNotes.filter((n) => n.is_archived);
                  }

                  return updatedNotes;
                });
              }}
              onDelete={(id) => {
                // Remove the deleted note from the list
                setNotes((prevNotes) => prevNotes.filter((n) => n.id !== id));
              }}
            />
          ))}
        </div>
      )}
    </div>
  );
};

export default Notes;