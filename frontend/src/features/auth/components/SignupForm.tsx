import { useState } from 'react';
import { useAppDispatch, useAppSelector } from '@/app/hooks';
import { signupUser } from '../thunks';

const SignupForm = () => {
  const dispatch = useAppDispatch();
  const { loading, error } = useAppSelector((state) => state.auth);

  const [formData, setFormData] = useState({ name: '', email: '', password: '' });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    dispatch(signupUser(formData));
  };

  return (
    <form onSubmit={handleSubmit}>
      <input name="name" value={formData.name} onChange={handleChange} placeholder="Name" required />
      <input name="email" value={formData.email} onChange={handleChange} placeholder="Email" required />
      <input name="password" type="password" value={formData.password} onChange={handleChange} placeholder="Password" required />
      <button type="submit" disabled={loading}>Sign Up</button>
      {error && <p>{error}</p>}
    </form>
  );
};

export default SignupForm;