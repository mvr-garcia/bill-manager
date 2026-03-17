-- Create bill_audits table for audit logs
CREATE TABLE bill_audits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bill_id UUID NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    performed_by UUID NOT NULL,
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for audit queries
CREATE INDEX idx_bill_audits_bill_id ON bill_audits(bill_id);
CREATE INDEX idx_bill_audits_performed_by ON bill_audits(performed_by);
CREATE INDEX idx_bill_audits_action ON bill_audits(action);
CREATE INDEX idx_bill_audits_created_at ON bill_audits(created_at);
