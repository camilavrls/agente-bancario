package policy

func CanAccessCustomerProfile(user AuthenticatedUser, targetCustomerID string) PolicyDecision {

	switch user.Role {
	case RoleCustomer:
		if user.CustomerID == targetCustomerID {
			return PolicyDecision{
				Allowed: true,
				Reason:  "allowed_customer_own_profile",
			}
		}
		return PolicyDecision{
			Allowed: false,
			Reason:  "denied_customer_other_profile",
		}

	case RoleManager:
		return PolicyDecision{
			Allowed: true,
			Reason:  "allowed_manager_customer_profile",
		}

	case RoleAdmin:
		return PolicyDecision{
			Allowed: true,
			Reason:  "allowed_admin_customer_profile",
		}

	default:
		return PolicyDecision{
			Allowed: false,
			Reason:  "denied_unknown_role",
		}
	}

}
