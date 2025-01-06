import { useEffect, useState, useContext } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { activateAccount } from "../services/api";
import { AuthContext } from "../contexts/AuthContext";
import { Loading } from "../components/Loading";

export const AccountActivation = () => {
  const { token } = useParams();
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const { showMessage } = useContext(AuthContext);

  useEffect(() => {
    const activate = async () => {
      try {
        setLoading(true);
        const response = await activateAccount(token);
        showMessage(
          "Account activated successfully! You can now log in.",
          "success",
        );
        navigate("/login");
      } catch (error) {
        showMessage(error.message || "Failed to activate account", "error");
        navigate("/");
      } finally {
        setLoading(false);
      }
    };

    if (token) {
      activate();
    }
  }, [token, navigate, showMessage]);

  if (loading) {
    return <Loading />;
  }

  return null;
};
