import { motion } from 'motion/react';
import { Card } from '../../../components/common';

function NotificationCard({ notifications }) {
  notifications = notifications.slice().reverse().slice(0, 3);
  if (notifications.length === 3) {
    notifications[2] = { title: '...', type: 'info' };
  }
  return (
    <Card variant="default" padding="md" animate className="overflow-hidden">
      <div className="flex items-center justify-between mb-3">
        <span className="text-neutral-400 text-sm">Notifications</span>
      </div>

      <div className="space-y-2">
        {notifications.map((notification, index) => (
          <motion.div
            key={index}
            className="text-sm"
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: index * 0.1 }}
          >
            <span
              className={`
                          ${
                            notification.type === 'success'
                              ? 'text-green-400'
                              : notification.type === 'info'
                                ? 'text-geek-400'
                                : 'text-yellow-400'
                          }
                      `}
            >
              •
            </span>
            <span className="text-neutral-300 ml-2">{notification.title}</span>
          </motion.div>
        ))}
      </div>
    </Card>
  );
}

export default NotificationCard;
