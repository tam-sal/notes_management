import { FaLinkedin, FaGithub, FaGlobe } from 'react-icons/fa';

const Footer = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="shadow-lg z-40 fixed bottom-0 w-full bg-slate-500">
      <div className="flex justify-center items-center py-2 px-4">
        <p className="text-white text-sm mr-8">
          &copy; {currentYear} NoteIT. All rights reserved.
        </p>
        <div className="flex gap-6">
          <a
            href="https://www.linkedin.com/in/tamsaleh/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-white text-xl"
          >
            <FaLinkedin />
          </a>
          <a
            href="https://github.com/tam-sal"
            target="_blank"
            rel="noopener noreferrer"
            className="text-white text-xl"
          >
            <FaGithub />
          </a>
          <a
            href="https://tamers-dev.vercel.app/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-white text-xl"
          >
            <FaGlobe />
          </a>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
