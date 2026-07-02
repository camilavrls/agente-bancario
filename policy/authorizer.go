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

func CanUpdateCardLimit(user AuthenticatedUser, targetCustomerID string) PolicyDecision {

	switch user.Role {
	case RoleCustomer:
		if user.CustomerID == targetCustomerID {
			return PolicyDecision{
				Allowed: true,
				Reason:  "allowed_customer_update_own_card_limit",
			}
		}
		return PolicyDecision{
			Allowed: false,
			Reason:  "denied_customer_update_other_card_limit",
		}

	case RoleManager:
		return PolicyDecision{
			Allowed: true,
			Reason:  "allowed_manager_update_card_limit",
		}

	case RoleAdmin:
		return PolicyDecision{
			Allowed: true,
			Reason:  "allowed_admin_update_card_limit",
		}

	default:
		return PolicyDecision{
			Allowed: false,
			Reason:  "denied_unknown_role",
		}
	}

}

func CanCreatePix(user AuthenticatedUser, fromCustomerID string) PolicyDecision {

	switch user.Role {
	case RoleCustomer:
		if user.CustomerID == fromCustomerID {
			return PolicyDecision{
				Allowed: true,
				Reason:  "allowed_customer_create_own_pix",
			}
		}
		return PolicyDecision{
			Allowed: false,
			Reason:  "denied_customer_create_other_pix",
		}

	case RoleManager:
		return PolicyDecision{
			Allowed: false,
			Reason:  "denied_manager_create_pix",
		}

	case RoleAdmin:
		return PolicyDecision{
			Allowed: false,
			Reason:  "denied_admin_create_pix",
		}

	default:
		return PolicyDecision{
			Allowed: false,
			Reason:  "denied_unknown_role",
		}
	}

}
