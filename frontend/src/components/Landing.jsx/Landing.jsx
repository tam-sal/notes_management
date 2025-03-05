import bg from "../../assets/bg-brown-1.png"
import { Link } from "react-router-dom";

const HomePage = () => {
  return (
    <div
      className="min-h-[80vh] bg-cover bg-center flex items-center rounded-lg justify-center"
      style={{ backgroundImage: `url(${bg})` }}
    >
      <div className="text-center">
        <div className="txt bg-slate-400 bg-opacity-50 rounded-xl">
          <h1 className="text-5xl font-bold ">
            Take Notes, Organize Life, Thrive Mentally
          </h1>
          <p className="text-2xl mb-8 ">
            Capture ideas, stay organized, and boost your mental and psychological well-being.
          </p>
        </div>

        <Link to="/register" className="btn border-t-orange-500 shadow-text-shadow">
          Get Started
        </Link>
      </div>

    </div>
  );
};

export default HomePage;