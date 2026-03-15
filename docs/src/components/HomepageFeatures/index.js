import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: '多题型支持',
    icon: '🎯',
    description: (
      <>
        支持问答题、静态题、动态附件题、容器题四种类型。
        每道题可配置多个 flag，支持静态、动态、UUID 三种 flag 格式。
      </>
    ),
  },
  {
    title: '动态容器靶机',
    icon: '🐳',
    description: (
      <>
        基于 Kubernetes 为每支队伍部署独立靶机，支持 Pod 网络和
        KubeOVN VPC 网络隔离，支持流量捕获和 FRP 端口转发。
      </>
    ),
  },
  {
    title: '完整赛事管理',
    icon: '🏆',
    description: (
      <>
        多种计分策略（静态/线性/对数）、三血奖励、RBAC 权限控制、
        内置作弊检测、Webhook 集成、Prometheus 监控。
      </>
    ),
  },
];

function Feature({icon, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center" style={{fontSize: '4rem', lineHeight: '1'}}>
        {icon}
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
