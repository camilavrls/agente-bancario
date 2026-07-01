package policy

func CanAccessCustomerResource(user AuthenticatedUser, targetCustomerID string) PolicyDecision {

	switch user.Role {
	case RoleCustomer:
		if user.CustomerID == targetCustomerID {
			return PolicyDecision{
				Allowed: true,
				Reason:  "allowed_customer_own_resource",
			}
		}
		return PolicyDecision{
			Allowed: false,
			Reason:  "denied_customer_other_resource",
		}

	case RoleManager:
		return PolicyDecision{
			Allowed: true,
			Reason:  "allowed_manager_customer_resource",
		}

	case RoleAdmin:
		return PolicyDecision{
			Allowed: true,
			Reason:  "allowed_admin_customer_resource",
		}

	default:
		return PolicyDecision{
			Allowed: false,
			Reason:  "denied_unknown_role",
		}
	}

}
