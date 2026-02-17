import { motion } from 'motion/react';
import { Card } from '../../../components/common';
import { useTranslation } from 'react-i18next';

function TechStack({
  technologies = [
    {
      category: 'Frontend',
      items: [
        { name: 'React', icon: '⚛️', description: 'UI Library' },
        { name: 'TailwindCSS', icon: '🎨', description: 'Styling Framework' },
        { name: 'Motion', icon: '✨', description: 'Animation Library' },
      ],
    },
    {
      category: 'Backend',
      items: [
        { name: 'Node.js', icon: '🟢', description: 'Runtime Environment' },
        { name: 'PostgreSQL', icon: '🐘', description: 'Database' },
        { name: 'Redis', icon: '🔄', description: 'Cache Layer' },
      ],
    },
    {
      category: 'Infrastructure',
      items: [
        { name: 'Docker', icon: '🐳', description: 'Containerization' },
        { name: 'Kubernetes', icon: '⚓️', description: 'Orchestration' },
        { name: 'AWS', icon: '☁️', description: 'Cloud Platform' },
      ],
    },
  ],
  developers = [
    {
      name: 'Alex Chen',
      role: 'Frontend Lead',
      picture: 'https://avatars.githubusercontent.com/u/1',
      github: 'https://github.com/alex',
      twitter: 'https://twitter.com/alex',
    },
    {
      name: 'Sarah Kim',
      role: 'Backend Lead',
      picture: 'https://avatars.githubusercontent.com/u/2',
      github: 'https://github.com/sarah',
      twitter: 'https://twitter.com/sarah',
    },
  ],
}) {
  const { t } = useTranslation();
  return (
    <div className="w-full max-w-[1000px] mx-auto space-y-12">
      {/* 技术栈部分 */}
      <motion.div className="space-y-8" initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        <h2 className="text-2xl text-neutral-50 font-mono border-b border-neutral-300/30 pb-4">
          {t('support.techStack.title')}
        </h2>
        <div className="grid grid-cols-3 gap-6">
          {technologies.map((tech, index) => (
            <motion.div
              key={index}
              className="space-y-4"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
            >
              <h3 className="text-neutral-50 font-mono">{tech.category}</h3>
              <div className="space-y-2">
                {tech.items.map((item, itemIndex) => (
                  <motion.div key={itemIndex} whileHover={{ x: 5 }}>
                    <Card
                      variant="default"
                      padding="sm"
                      className="hover:border-geek-400 hover:shadow-[0_0_20px_rgba(89,126,247,0.4)]
                                            transition-all duration-200"
                    >
                      <div className="flex items-center gap-3">
                        <span className="text-xl">{item.icon}</span>
                        <div>
                          <div className="text-neutral-50 font-mono">{item.name}</div>
                          <div className="text-neutral-400 text-sm">{item.description}</div>
                        </div>
                      </div>
                    </Card>
                  </motion.div>
                ))}
              </div>
            </motion.div>
          ))}
        </div>
      </motion.div>

      {/* 开发者部分 */}
      <motion.div
        className="space-y-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
      >
        <h2 className="text-2xl text-neutral-50 font-mono border-b border-neutral-300/30 pb-4">
          {t('support.team.title')}
        </h2>
        <div className="grid grid-cols-2 gap-6">
          {developers.map((dev, index) => (
            <motion.div key={index} whileHover={{ y: -2 }}>
              <Card variant="default" padding="md" className="hover:border-geek-400 transition-all duration-200">
                <div className="flex items-center gap-4">
                  <div className="w-16 h-16 rounded-md border-2 border-neutral-300 overflow-hidden">
                    <img src={dev.picture} alt={dev.name} className="w-full h-full object-cover" />
                  </div>
                  <div>
                    <h3 className="text-neutral-50 font-mono text-lg">{dev.name}</h3>
                    <p className="text-geek-400 text-sm">{dev.role}</p>
                    <div className="flex gap-3 mt-2">
                      <a
                        href={dev.github}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-neutral-400 hover:text-neutral-50 transition-colors"
                      >
                        {t('common.github')}
                      </a>
                    </div>
                  </div>
                </div>
              </Card>
            </motion.div>
          ))}
        </div>
      </motion.div>
    </div>
  );
}

export default TechStack;
