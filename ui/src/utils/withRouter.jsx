import { useNavigate, useParams, useLocation } from 'react-router-dom';

/**
 * Replacement for react-router-dom v5's withRouter HOC.
 * Provides navigate, params, and location as props to class components.
 */
export function withRouter(Component) {
  function ComponentWithRouterProps(props) {
    const navigate = useNavigate();
    const params = useParams();
    const location = useLocation();
    
    return (
      <Component
        {...props}
        navigate={navigate}
        params={params}
        location={location}
      />
    );
  }
  
  return ComponentWithRouterProps;
}
