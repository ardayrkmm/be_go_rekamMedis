"""
Script Export Data Firestore -> MySQL
======================================
Cara pakai:
1. pip install firebase-admin mysql-connector-python
2. python scripts/export_firestore_to_mysql.py
"""

import firebase_admin
from firebase_admin import credentials, firestore
import mysql.connector
from datetime import datetime
import os
import time

FIREBASE_CRED_PATH = os.path.join(os.path.dirname(__file__), '..', 'rekamGo.json')

MYSQL_CONFIG = {
    'host': 'localhost',
    'port': 3306,
    'user': 'root',
    'password': '',
    'database': 'rekam_medis',
    'charset': 'utf8mb4',
}


def stream_with_retry(collection_ref, max_retries=5):
    """Stream Firestore collection dengan retry otomatis kalau kena 429."""
    for attempt in range(max_retries):
        try:
            return list(collection_ref.stream())
        except Exception as e:
            if '429' in str(e) or 'Quota' in str(e) or 'RESOURCE_EXHAUSTED' in str(e):
                wait = (2 ** attempt) * 10  # 10, 20, 40, 80, 160 detik
                print(f"  [429 Quota] Tunggu {wait} detik lalu coba lagi (attempt {attempt+1}/{max_retries})...")
                time.sleep(wait)
            else:
                raise
    raise Exception(f"Gagal setelah {max_retries} percobaan")


def init_firebase():
    cred = credentials.Certificate(FIREBASE_CRED_PATH)
    firebase_admin.initialize_app(cred)
    return firestore.client()


def ts(val):
    if val is None:
        return datetime.now().isoformat()
    if hasattr(val, 'isoformat'):
        return val.isoformat()
    return str(val)


def ts_null(val):
    if val is None:
        return None
    if hasattr(val, 'isoformat'):
        return val.isoformat()
    return str(val)


def s(val, max_len=255):
    if val is None:
        return None
    return str(val)[:max_len]


def export_users(db, cursor):
    print("\n[1/10] Exporting users...")
    count = 0
    for doc in db.collection('users').stream():
        d = doc.to_dict()
        try:
            cursor.execute("""
                INSERT IGNORE INTO users (id, name, email, password, role, photo, created_at, updated_at)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
            """, (
                doc.id, s(d.get('name')), s(d.get('email')),
                s(d.get('password'), 500), s(d.get('role')),
                s(d.get('photo'), 500),
                ts(d.get('createdAt')), ts(d.get('updatedAt')),
            ))
            count += 1
        except Exception as e:
            print(f"  Skip user {doc.id}: {e}")
    print(f"  OK {count} users")


def export_patient_categories(db, cursor):
    print("\n[2/10] Exporting patient_categories...")
    count = 0
    for doc in db.collection('patientCategories').stream():
        d = doc.to_dict()
        try:
            cursor.execute("""
                INSERT IGNORE INTO patient_categories (id, name, created_at, updated_at)
                VALUES (%s, %s, %s, %s)
            """, (
                doc.id, s(d.get('name')),
                ts(d.get('createdAt')), ts(d.get('updatedAt')),
            ))
            count += 1
        except Exception as e:
            print(f"  Skip {doc.id}: {e}")
    print(f"  OK {count} patient_categories")


def export_patients(db, cursor):
    print("\n[3/10] Exporting patients...")
    count = 0
    for doc in db.collection('patients').stream():
        d = doc.to_dict()
        try:
            cursor.execute("""
                INSERT IGNORE INTO patients (
                    id, medical_record_number, nik, name, birth_date,
                    patient_category_id, gender_id, blood_type, address, phone, email,
                    occupation, marital_status, emergency_contact_name, emergency_contact_phone,
                    medical_history, allergies, created_at, updated_at
                ) VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)
            """, (
                doc.id, s(d.get('medicalRecordNumber')), s(d.get('nik')),
                s(d.get('name')), ts_null(d.get('birthDate')),
                s(d.get('patientCategoryId')), s(d.get('genderId')),
                s(d.get('bloodType')), s(d.get('address'), 500),
                s(d.get('phone')), s(d.get('email')),
                s(d.get('occupation')), s(d.get('maritalStatus')),
                s(d.get('emergencyContactName')), s(d.get('emergencyContactPhone')),
                s(d.get('medicalHistory'), 2000), s(d.get('allergies'), 1000),
                ts(d.get('createdAt')), ts(d.get('updatedAt')),
            ))
            count += 1
        except Exception as e:
            print(f"  Skip patient {doc.id}: {e}")
    print(f"  OK {count} patients")


def export_physiotherapists(db, cursor):
    print("\n[4/10] Exporting physiotherapists...")
    count = 0
    for doc in db.collection('physiotherapists').stream():
        d = doc.to_dict()
        try:
            cursor.execute("""
                INSERT IGNORE INTO physiotherapists (
                    id, name, specialization, sip, phone, email,
                    address, gender, photo, status, created_at, updated_at
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            """, (
                doc.id, s(d.get('name')), s(d.get('specialization')),
                s(d.get('sip')), s(d.get('phone')), s(d.get('email')),
                s(d.get('address'), 500), s(d.get('gender')),
                s(d.get('photo'), 500), s(d.get('status')),
                ts(d.get('createdAt')), ts(d.get('updatedAt')),
            ))
            count += 1
        except Exception as e:
            print(f"  Skip {doc.id}: {e}")
    print(f"  OK {count} physiotherapists")


def export_service_masters(db, cursor):
    print("\n[5/10] Exporting service_masters...")
    count = 0
    for doc in db.collection('serviceMasters').stream():
        d = doc.to_dict()
        try:
            cursor.execute("""
                INSERT IGNORE INTO service_masters (id, name, description, price, duration, created_at, updated_at)
                VALUES (%s, %s, %s, %s, %s, %s, %s)
            """, (
                doc.id, s(d.get('name')), s(d.get('description'), 1000),
                d.get('price') or 0, d.get('duration') or 0,
                ts(d.get('createdAt')), ts(d.get('updatedAt')),
            ))
            count += 1
        except Exception as e:
            print(f"  Skip {doc.id}: {e}")
    print(f"  OK {count} service_masters")


def export_appointments(db, cursor):
    print("\n[6/10] Exporting appointments...")
    count = 0
    for doc in db.collection('appointments').stream():
        d = doc.to_dict()
        try:
            cursor.execute("""
                INSERT IGNORE INTO appointments (
                    id, visit_number, patient_id, physiotherapist_id, service_master_id,
                    appointment_date, appointment_time, complaint, status, notes,
                    created_at, updated_at
                ) VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)
            """, (
                doc.id, s(d.get('visitNumber')),
                s(d.get('patientId')), s(d.get('physiotherapistId')),
                s(d.get('serviceMasterId')), ts_null(d.get('appointmentDate')),
                s(d.get('appointmentTime')), s(d.get('complaint'), 1000),
                s(d.get('status')), s(d.get('notes'), 1000),
                ts(d.get('createdAt')), ts(d.get('updatedAt')),
            ))
            count += 1
        except Exception as e:
            print(f"  Skip {doc.id}: {e}")
    print(f"  OK {count} appointments")


def export_medical_records(db, cursor):
    print("\n[7/10] Exporting medical_records...")
    count = 0
    for doc in db.collection('medicalRecords').stream():
        d = doc.to_dict()
        try:
            cursor.execute("""
                INSERT IGNORE INTO medical_records (
                    id, visit_number, patient_id, service_id, physiotherapist_id,
                    appointment_id, examination_date, anamnesis, diagnosis, therapy,
                    notes, created_at, updated_at
                ) VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)
            """, (
                doc.id, s(d.get('visitNumber')),
                s(d.get('patientId')), s(d.get('serviceId')),
                s(d.get('physiotherapistId')), s(d.get('appointmentId')),
                ts_null(d.get('examinationDate')),
                s(d.get('anamnesis'), 2000), s(d.get('diagnosis'), 1000),
                s(d.get('therapy'), 1000), s(d.get('notes'), 1000),
                ts(d.get('createdAt')), ts(d.get('updatedAt')),
            ))
            count += 1
        except Exception as e:
            print(f"  Skip {doc.id}: {e}")
    print(f"  OK {count} medical_records")


def export_payments(db, cursor):
    print("\n[8/10] Exporting payments...")
    count = 0
    for doc in db.collection('payments').stream():
        d = doc.to_dict()
        try:
            cursor.execute("""
                INSERT IGNORE INTO payments (
                    id, invoice_number, appointment_id, therapy_session_id,
                    patient_id, patient_name, physiotherapist_id, physiotherapist_name,
                    payment_date, payment_method, status,
                    subtotal, discount, tax, total,
                    notes, created_at, updated_at
                ) VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)
            """, (
                doc.id, s(d.get('invoiceNumber')),
                s(d.get('appointmentId')), s(d.get('therapySessionId')),
                s(d.get('patientId')), s(d.get('patientName')),
                s(d.get('physiotherapistId')), s(d.get('physiotherapistName')),
                ts_null(d.get('paymentDate')), s(d.get('paymentMethod')),
                s(d.get('status')),
                d.get('subtotal') or d.get('totalAmount') or 0,
                d.get('discount') or 0,
                d.get('tax') or 0,
                d.get('total') or d.get('totalAmount') or 0,
                s(d.get('notes'), 1000),
                ts(d.get('createdAt')), ts(d.get('updatedAt')),
            ))
            count += 1
        except Exception as e:
            print(f"  Skip {doc.id}: {e}")
    print(f"  OK {count} payments")


def export_notifications(db, cursor):
    print("\n[9/10] Exporting notifications...")
    count = 0
    for doc in db.collection('notifications').stream():
        d = doc.to_dict()
        try:
            cursor.execute("""
                INSERT IGNORE INTO notifications (
                    id, user_id, title, message, type, is_read, created_at, updated_at
                ) VALUES (%s,%s,%s,%s,%s,%s,%s,%s)
            """, (
                doc.id, s(d.get('userId')), s(d.get('title')),
                s(d.get('message'), 1000), s(d.get('type')),
                bool(d.get('isRead', False)),
                ts(d.get('createdAt')), ts(d.get('updatedAt')),
            ))
            count += 1
        except Exception as e:
            print(f"  Skip {doc.id}: {e}")
    print(f"  OK {count} notifications")


def export_activity_logs(db, cursor):
    print("\n[10/10] Exporting activity_logs...")
    count = 0
    for doc in db.collection('activityLogs').stream():
        d = doc.to_dict()
        try:
            cursor.execute("""
                INSERT IGNORE INTO activity_logs (
                    id, user_id, action, entity_type, entity_id, properties, ip_address, created_at, updated_at
                ) VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s)
            """, (
                doc.id, s(d.get('userId')), s(d.get('action')),
                s(d.get('entityType')), s(d.get('entityId')),
                s(d.get('properties') or d.get('description'), 1000),
                s(d.get('ipAddress')),
                ts(d.get('createdAt')), ts(d.get('updatedAt')),
            ))
            count += 1
        except Exception as e:
            print(f"  Skip {doc.id}: {e}")
    print(f"  OK {count} activity_logs")


def main():
    print("=" * 50)
    print("  Export Firestore -> MySQL")
    print("=" * 50)

    print("\nMenghubungkan ke Firestore...")
    db = init_firebase()
    print("  OK Firestore terhubung")

    print("\nMenghubungkan ke MySQL...")
    conn = mysql.connector.connect(**MYSQL_CONFIG)
    conn.autocommit = False
    cursor = conn.cursor()
    print("  OK MySQL terhubung")

    try:
        cursor.execute("SET FOREIGN_KEY_CHECKS=0")

        export_users(db, cursor)
        time.sleep(3)
        export_patient_categories(db, cursor)
        time.sleep(3)
        export_patients(db, cursor)
        time.sleep(3)
        export_physiotherapists(db, cursor)
        time.sleep(3)
        export_service_masters(db, cursor)
        time.sleep(3)
        export_appointments(db, cursor)
        time.sleep(3)
        export_medical_records(db, cursor)
        time.sleep(3)
        export_payments(db, cursor)
        time.sleep(3)
        export_notifications(db, cursor)
        time.sleep(3)
        export_activity_logs(db, cursor)

        cursor.execute("SET FOREIGN_KEY_CHECKS=1")
        conn.commit()
        print("\n\nSemua data berhasil di-export ke MySQL!")

    except Exception as e:
        conn.rollback()
        print(f"\nError: {e}")
        raise
    finally:
        cursor.close()
        conn.close()


if __name__ == '__main__':
    main()
