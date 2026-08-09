import os
import re

model_dir = r'e:\Bisnis\RekamMedis\backend_go_firebase\internal\models'
for filename in os.listdir(model_dir):
    if not filename.endswith('.go'):
        continue
    path = os.path.join(model_dir, filename)
    with open(path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Replace gorm:"-" with gorm:"foreignKey:X;constraint:-"
    replacements = [
        (r'Category\s+\*PatientCategory\s+`json:\"category,omitempty\"\s+gorm:\"-\"`', r'Category *PatientCategory `json:"category,omitempty" gorm:"foreignKey:PatientCategoryID;constraint:-"`'),
        (r'GenderData\s+\*Gender\s+`json:\"gender_data,omitempty\"\s+gorm:\"-\"`', r'GenderData *Gender `json:"gender_data,omitempty" gorm:"foreignKey:GenderID;constraint:-"`'),
        (r'Appointments\s+\[\]Appointment\s+`json:\"appointments,omitempty\"\s+gorm:\"-\"`', r'Appointments []Appointment `json:"appointments,omitempty" gorm:"foreignKey:PatientID;constraint:-"`'),
        (r'MedicalRecords\s+\[\]MedicalRecord\s+`json:\"medical_records,omitempty\"\s+gorm:\"-\"`', r'MedicalRecords []MedicalRecord `json:"medical_records,omitempty" gorm:"foreignKey:PatientID;constraint:-"`'),
        
        (r'Patient\s+\*Patient\s+`json:\"patient,omitempty\"\s+gorm:\"-\"`', r'Patient *Patient `json:"patient,omitempty" gorm:"foreignKey:PatientID;constraint:-"`'),
        (r'Physiotherapist\s+\*Physiotherapist\s+`json:\"physiotherapist,omitempty\"\s+gorm:\"-\"`', r'Physiotherapist *Physiotherapist `json:"physiotherapist,omitempty" gorm:"foreignKey:PhysiotherapistID;constraint:-"`'),
        (r'ServiceMaster\s+\*ServiceMaster\s+`json:\"service_master,omitempty\"\s+gorm:\"-\"`', r'ServiceMaster *ServiceMaster `json:"service_master,omitempty" gorm:"foreignKey:ServiceID;constraint:-"`'),
        (r'Service\s+\*ServiceMaster\s+`json:\"service,omitempty\"\s+gorm:\"-\"`', r'Service *ServiceMaster `json:"service,omitempty" gorm:"foreignKey:ServiceID;constraint:-"`'),
        
        (r'TherapySession\s+\*TherapySession\s+`json:\"therapy_session,omitempty\"\s+gorm:\"-\"`', r'TherapySession *TherapySession `json:"therapy_session,omitempty" gorm:"foreignKey:AppointmentID;constraint:-"`'),
        (r'TherapySession\s+\*TherapySession\s+`gorm:\"-\"\s+json:\"therapy_session,omitempty\"`', r'TherapySession *TherapySession `gorm:"foreignKey:TherapySessionID;constraint:-" json:"therapy_session,omitempty"`'),
        
        (r'Appointment\s+\*Appointment\s+`json:\"appointment,omitempty\"\s+gorm:\"-\"`', r'Appointment *Appointment `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID;constraint:-"`'),
        (r'PaymentDetails\s+\[\]PaymentDetail\s+`json:\"payment_details,omitempty\"\s+gorm:\"-\"`', r'PaymentDetails []PaymentDetail `json:"payment_details,omitempty" gorm:"foreignKey:PaymentID;constraint:-"`'),
        (r'Payment\s+\*Payment\s+`json:\"payment,omitempty\"\s+gorm:\"-\"`', r'Payment *Payment `json:"payment,omitempty" gorm:"foreignKey:PaymentID;constraint:-"`'),
        (r'User\s+\*User\s+`json:\"user,omitempty\"\s+gorm:\"-\"`', r'User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:-"`'),
    ]
    
    new_content = content
    for old, new in replacements:
        new_content = re.sub(old, new, new_content)
        
    if new_content != content:
        with open(path, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f'Updated {filename}')
