// src/features/auth/components/SignupForm.tsx
import { useFormik } from "formik";
import * as Yup from "yup";
import { TextInput, PasswordInput, Button, Select } from "@mantine/core";
import { useAppDispatch } from "@/app/hooks";
import { signupUser } from "../../thunks";

const SignupForm = () => {
  const dispatch = useAppDispatch();

  const formik = useFormik({
    initialValues: {
      name: "",
      email: "",
      password: "",
      userType: "individual",
    },
    validationSchema: Yup.object({
      name: Yup.string().required("Required"),
      email: Yup.string().email("Invalid email").required("Required"),
      password: Yup.string().min(6).required("Required"),
      userType: Yup.string().oneOf(["individual", "team"]).required("Required"),
    }),
    onSubmit: (values) => {
      dispatch(signupUser(values));
    },
  });

  return (
    <form onSubmit={formik.handleSubmit}>
      <TextInput
        label="Name"
        name="name"
        value={formik.values.name}
        onChange={formik.handleChange}
        error={formik.touched.name && formik.errors.name}
      />
      <TextInput
        label="Email"
        name="email"
        value={formik.values.email}
        onChange={formik.handleChange}
        error={formik.touched.email && formik.errors.email}
        mt="sm"
      />
      <PasswordInput
        label="Password"
        name="password"
        value={formik.values.password}
        onChange={formik.handleChange}
        error={formik.touched.password && formik.errors.password}
        mt="sm"
      />
      <Select
        label="User Type"
        name="userType"
        data={[
          { value: "individual", label: "Individual" },
          { value: "team", label: "Team" },
        ]}
        value={formik.values.userType}
        onChange={(value) => formik.setFieldValue("userType", value)}
        error={formik.touched.userType && formik.errors.userType}
        mt="sm"
      />
      <Button fullWidth type="submit" mt="md">
        Sign Up
      </Button>
    </form>
  );
};

export default SignupForm;
