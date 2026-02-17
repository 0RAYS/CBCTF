import { Button, LanguageSwitcher, Avatar } from '../../common';

function NavBar({
  tabs = [],
  activeTab,
  onTabChange,
  logo = 'DEEP DIVE',
  pictureSrc = '',
  userName = '',
  className = '',
  onLogoClick = () => {},
  onPictureClick = () => {},
}) {
  return (
    <div
      className={`fixed top-0 left-0 w-full h-[80px] flex items-center justify-between px-12 bg-black/30 backdrop-blur-[2px] border-b border-neutral-300 ${className}`}
    >
      {/* 左侧区域包装 */}
      <div className="flex items-center">
        {/* Logo区域 */}
        <div className="relative">
          <Button
            variant="outline"
            size="lg"
            className="min-w-[180px] h-[50px] font-mono text-lg tracking-wider"
            onClick={onLogoClick}
          >
            {logo}
          </Button>
          {/* Logo装饰角 */}
          <div className="absolute -top-[3px] -right-[3px] w-[10px] h-[10px] border-t border-r border-neutral-300" />
          <div className="absolute -bottom-[3px] -left-[3px] w-[10px] h-[10px] border-b border-l border-neutral-300" />
        </div>

        {/* 标签区域 */}
        <div className="flex ml-16 gap-6">
          {tabs.map((tab) => (
            <Button
              key={tab.id}
              variant={activeTab === tab.id ? 'primary' : 'outline'}
              size="sm"
              className="min-w-[100px]"
              onClick={() => onTabChange(tab.id)}
            >
              {tab.label}
            </Button>
          ))}
        </div>
      </div>

      {/* 右侧语言切换与头像区域 */}
      <div className="flex items-center gap-4">
        <LanguageSwitcher size="sm" />
        <div className="relative group" onClick={onPictureClick}>
          <div className="w-[45px] h-[45px] border border-neutral-300 rounded-md overflow-hidden transition-all duration-200 hover:border-neutral-100 hover:shadow-[0_0_15px_rgba(179,179,179,0.3)] cursor-pointer">
            <Avatar src={pictureSrc} name={userName} size={45} shape="rounded" />
          </div>
          {/* 装饰角 */}
          <div className="absolute -top-[2px] -right-[2px] w-[8px] h-[8px] border-t border-r border-neutral-300 group-hover:border-neutral-100 rounded-tr-none"></div>
          <div className="absolute -bottom-[2px] -left-[2px] w-[8px] h-[8px] border-b border-l border-neutral-300 group-hover:border-neutral-100 rounded-bl-none"></div>
        </div>
      </div>
    </div>
  );
}

export default NavBar;
